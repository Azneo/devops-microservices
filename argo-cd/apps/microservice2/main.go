package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	_ "github.com/lib/pq"
)

func main() {
	// Kafka config
	broker := "kafka.microservices.svc.cluster.local:9092"
	topic := "microservice-topic"

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

	// Connect to PostgreSQL
	connStr := "postgres://user:pass@postgres-postgresql.microservices.svc.cluster.local:5432/messages?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS kafka_messages (
		id SERIAL PRIMARY KEY,
		content TEXT UNIQUE
	)`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Kafka consumer setup
	master, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log.Fatal("Kafka consumer error:", err)
	}
	defer master.Close()

	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Partition error:", err)
	}
	defer consumer.Close()

	fmt.Println("Kafka consumer started. Waiting for messages...")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg := <-consumer.Messages():
			message := string(msg.Value)
			fmt.Printf("Received: %s\n", message)

			// Insert only if not already present
			_, err := db.Exec("INSERT INTO kafka_messages (content) VALUES ($1) ON CONFLICT DO NOTHING", message)
			if err != nil {
				log.Println("Insert error:", err)
			}

		case err := <-consumer.Errors():
			log.Println("Kafka error:", err)
		case <-signals:
			fmt.Println("Interrupt received, shutting down.")
			return
		}
	}
}
