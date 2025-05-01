package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	kafkaBrokers := []string{"kafka.microservices.svc.cluster.local:9092"}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	config.Net.SASL.Enable = true
	config.Net.SASL.User = "user1"
	config.Net.SASL.Password = "C6lH3xSsg9"
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.Handshake = true
	config.Net.SASL.Version = sarama.SASLHandshakeV1

	config.Net.TLS.Enable = false
	config.Metadata.Full = true

	producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
	if err != nil {
		log.Fatal("Error creating Kafka producer:", err)
	}
	defer producer.Close()

	fmt.Println("Kafka producer started. Sending messages every 5 seconds...")

	for {
		message := &sarama.ProducerMessage{
			Topic: "microservice-topic",
			Value: sarama.StringEncoder("Hello from Microservice 1!"),
		}

		partition, offset, err := producer.SendMessage(message)
		if err != nil {
			log.Println("Error sending message to Kafka:", err)
		} else {
			fmt.Printf("Sent message to partition %d with offset %d\n", partition, offset)
		}

		time.Sleep(5 * time.Second)
	}
}
