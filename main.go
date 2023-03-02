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
		"isolation.level":"read_committed",
	})

	topic := "transactions"
	
	err = consumer.Assign([]kafka.TopicPartition{{Topic: &topic, Partition: 1,Offset: kafka.OffsetStored}})
	if err != nil {
		log.Fatal(err)
	}
	
	// low,_,err := consumer.QueryWatermarkOffsets(topic,2,2000)


	run := true
	

	for run {
		msg, err := consumer.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			consumer.CommitMessage(msg)
		} else if err != nil {
		// The client will automatically try to recover from all errors.
		// Timeout is not considered an error because it is raised by
		// ReadMessage in absence of messages.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	consumer.Close()

}
