package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

func main() {
	broker := "kafka.microservices.svc.cluster.local:9092"
	topic := "microservice-topic" // ✅ match producer topic

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	config.Net.SASL.Enable = true
	config.Net.SASL.User = "user1"
	config.Net.SASL.Password = "C6lH3xSsg9"
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.Handshake = true
	config.Net.SASL.Version = sarama.SASLHandshakeV1

	config.Net.TLS.Enable = false
	config.Version = sarama.V2_1_0_0

	master, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log.Fatal("Error creating consumer:", err)
	}
	defer master.Close()

	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Error consuming partition:", err)
	}
	defer consumer.Close()

	fmt.Println("Kafka consumer started. Waiting for messages...")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg := <-consumer.Messages():
			fmt.Printf("Received: %s\n", string(msg.Value))
		case err := <-consumer.Errors():
			log.Println("Error:", err)
		case <-signals:
			fmt.Println("Interrupt detected. Exiting...")
			return
		}
	}
}
