package main

import (
    "fmt"
    "log"
    "github.com/IBM/sarama"
)

func main() {
    // Kafka broker address
    kafkaBrokers := []string{"kafka.microservices.svc.cluster.local:9092"} // Use the Kafka service DNS in Kubernetes

    // Kafka Producer config
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
    if err != nil {
        log.Fatal("Error creating Kafka producer: ", err)
    }
    defer producer.Close()

    // Sending a message to Kafka
    message := &sarama.ProducerMessage{
        Topic: "microservice-topic", // Topic name
        Value: sarama.StringEncoder("Hello from Microservice 1!"),
    }

    // Send the message
    partition, offset, err := producer.SendMessage(message)
    if err != nil {
        log.Fatal("Error sending message to Kafka: ", err)
    }
    fmt.Printf("Message sent to partition %d with offset %d\n", partition, offset)
}
