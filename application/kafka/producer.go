package kafka

import (
	"log"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

func NewKafkaProducer() *ckafka.Producer {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
	}
	p, err := ckafka.NewProducer(configMap)
	if err != nil {
		panic(err)
	}
	return p
}

func Publish(msg string, topic string, producer *ckafka.Producer, deliveryChannel chan ckafka.Event) error {
	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Value: []byte(msg),
	}
	err := producer.Produce(message, deliveryChannel)
	if err != nil {
		return err
	}
	return nil
}

func DeliveryReport(deliveryChannel chan ckafka.Event) {
	for e := range deliveryChannel {
		switch ev := e.(type) {
		case *ckafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Print("Delivery failed: ", ev.TopicPartition)
			} else {
				log.Print("Delivered message to: ", ev.TopicPartition)
			}
		}
	}
}
