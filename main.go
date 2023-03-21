package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/loyalty-application/go-worker-node/models"
	"github.com/loyalty-application/go-worker-node/config"
	"github.com/loyalty-application/go-worker-node/services"
	"github.com/loyalty-application/go-worker-node/collections"
)

func main() {

	// Connect to DB
	config.DBinstance()

	// Create a new Kafka Conumer
	consumer, err := getKafkaConsumer()
	defer consumer.Close()

	// Subscribe to kafka topic
	topic := "ftptransactions"
	consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("start consuming ... !!")

	for {
		var transactions models.TransactionList

		for i := 0; i < 20000; i++ {
			msg, err := consumer.ReadMessage(time.Second)
			
			if err != nil { break }

			var transaction models.Transaction
			json.Unmarshal(msg.Value, &transaction)
			services.ConvertPoints(&transaction)
			// fmt.Println(transaction)  // DEBUG
			transactions.Transactions = append(transactions.Transactions, transaction)
		}

		// If there are transactions, insert them into the DB and commit
		if len(transactions.Transactions) != 0 {
			collections.CreateTransactions(transactions)
			consumer.Commit()
		}
	}
}


func getKafkaConsumer() (consumer *kafka.Consumer, err error) {
	server := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        server,
		"group.id":                 "FtpWorkerGroup",
		"client.id":                "FtpProcessing",
		"enable.auto.commit":       false,
		"enable.auto.offset.store": true,
		"auto.offset.reset":        "earliest",
		"isolation.level":          "read_committed",
	})
}