package cmd

import (
	"log"
	"os"

	"github.com/RuanScherer/codepix-go/application/grpc"
	"github.com/RuanScherer/codepix-go/application/kafka"
	"github.com/RuanScherer/codepix-go/infrastructure/db"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

var grpcPortNumber int

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run gRPC and Kafka Consumer",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("error loading .env file")
		}

		database := db.ConnectDB(os.Getenv("env"))
		go grpc.StartGRPCServer(database, portNumber)

		deliveryChannel := make(chan ckafka.Event)
		producer := kafka.NewKafkaProducer()

		go kafka.DeliveryReport(deliveryChannel)

		kafkaProcessor := kafka.NewKafkaProcessor(database, producer, deliveryChannel)
		kafkaProcessor.Consume()
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().IntVarP(&grpcPortNumber, "grpc-port", "p", 50051, "gRPC Server Port")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
