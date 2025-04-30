package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"

    "github.com/IBM/sarama"
)

func main() {
    // Kafka broker address
    broker := "kafka.microservices.svc.cluster.local:9092"
    topic := "example-topic"

    config := sarama.NewConfig()
    config.Consumer.Return.Errors = true

    master, err := sarama.NewConsumer([]string{broker}, config)
    if err != nil {
        panic(err)
    }
    defer master.Close()

    consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
    if err != nil {
        panic(err)
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
            fmt.Println("Interrupt is detected. Exiting...")
            return
        }
    }
}
