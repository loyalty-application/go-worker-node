package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "FtpWorkerGroup",
		"client.id": "producer2",
		"enable.auto.commit":true,
		"enable.auto.offset.store": false,
		"auto.commit.interval.ms": "4000",
	})

	topic := "transaction"

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	run := true

	for run {
		msg, err := consumer.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		}
		// else if !err.(kafka.Error).IsTimeout() {
		//// The client will automatically try to recover from all errors.
		//// Timeout is not considered an error because it is raised by
		//// ReadMessage in absence of messages.
		//fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		//}
	}

	consumer.Close()

}
