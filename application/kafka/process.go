package kafka

import (
	"log"
	"os"

	"github.com/RuanScherer/codepix-go/application/factory"
	appmodel "github.com/RuanScherer/codepix-go/application/model"
	"github.com/RuanScherer/codepix-go/domain/model"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

type KafkaProcessor struct {
	database        *gorm.DB
	producer        *ckafka.Producer
	deliveryChannel chan ckafka.Event
}

func NewKafkaProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChannel chan ckafka.Event) *KafkaProcessor {
	return &KafkaProcessor{
		database:        database,
		producer:        producer,
		deliveryChannel: deliveryChannel,
	}
}

func (k *KafkaProcessor) Consume() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	}
	c, err := ckafka.NewConsumer(configMap)
	if err != nil {
		panic(err)
	}

	topics := []string{os.Getenv("kafkaTransactionTopic"), os.Getenv("kafkaTransactionConfirmationTopic")}
	c.SubscribeTopics(topics, nil)
	log.Print("Kafka consumer has been started")

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			k.processMessage(msg)
		}
	}
}

func (k *KafkaProcessor) processMessage(msg *ckafka.Message) {
	transactionsTopic := os.Getenv("kafkaTransactionTopic")
	transactionConfirmationTopic := os.Getenv("kafkaTransactionConfirmationTopic")

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		k.processTransaction(msg)
	case transactionConfirmationTopic:
		k.processTransactionConfirmation(msg)
	default:
		log.Print("not a valid topic", string(msg.Value))
	}
}

func (k *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJSON(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(k.database)
	createdTransaction, err := transactionUseCase.Register(
		transaction.AccountID,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyKindTo,
		transaction.Description,
	)
	if err != nil {
		log.Println("error registering transaction", err)
		return err
	}

	targetBankTopic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending

	transactionJSON, err := transaction.ToJSON()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJSON), targetBankTopic, k.producer, k.deliveryChannel)
	return err
}

func (k *KafkaProcessor) processTransactionConfirmation(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJSON(msg.Value)
	if err != nil {
		return err
	}

	if transaction.Status == model.TransactionConfirmed {
		err = k.confirmTransaction(transaction)
		if err != nil {
			return err
		}
	} else if transaction.Status == model.TransactionCompleted {
		err = k.completeTransaction(transaction)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *KafkaProcessor) confirmTransaction(transaction *appmodel.Transaction) error {
	transactionUseCase := factory.TransactionUseCaseFactory(k.database)
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)
	if err != nil {
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	transactionJSON, err := transaction.ToJSON()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJSON), topic, k.producer, k.deliveryChannel)
	return err
}

func (k *KafkaProcessor) completeTransaction(transaction *appmodel.Transaction) error {
	transactionUseCase := factory.TransactionUseCaseFactory(k.database)
	_, err := transactionUseCase.Complete(transaction.ID)
	return err
}
