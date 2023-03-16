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
			// TODO: Process transaction
			
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			// Only commit after successfully processed the message
			consumer.CommitMessage(msg)
		} else if err != nil {
			// TODO Handle error
		}
	}

	consumer.Close()

}
