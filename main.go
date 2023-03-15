package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

/*
 * TODO:
 * 1. Migrate POST method from controller.transaction.go after Gabriel makes it ATOMIC
 * 2. Implement VALIDATION checks (TBC)
 * 3. Implement reading from Kafka
 */

func main() {

	// Setting up a connection with kafka
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        "localhost:9092",
		"group.id":                 "FtpWorkerGroup",
		"client.id":                "FtpProcessing",
		"enable.auto.commit":       false,
		"enable.auto.offset.store": false,
		"auto.offset.reset":        "earliest",
		"isolation.level":          "read_committed",
	})

	// Creating a topic categorization
	topic := "ftptransactions"

	// Subscribe to the message broker with decided topic
	consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	run := true

	// Run a infinite loop that constantly checks for messages
	for run {

		msg, err := consumer.ReadMessage(time.Second)

		if err == nil {
			// native, _, err := codec.NativeFromBinary(msg.Value)
			// if err != nil {
			// 	fmt.Println("Error decoding Avro message value:", err)
			// 	continue
			// }
			
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			// Only commit after successfully processed the message
			consumer.CommitMessage(msg)
		} else if err != nil {
			// The client will automatically try to recover from all errors.
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			// fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	consumer.Close()

}
