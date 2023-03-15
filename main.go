package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        "localhost:9092",
		"group.id":                 "FtpWorkerGroup",
		"client.id":                "FtpProcessing",
		"enable.auto.commit":       false,
		"enable.auto.offset.store": false,
		"auto.offset.reset":        "earliest",
		"isolation.level":          "read_committed",
	})

	topic := "ftptransactions"

	consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	run := true

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
