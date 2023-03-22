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

		// Process transactions in batches
		for i := 0; i < 50000; i++ {
			msg, err := consumer.ReadMessage(time.Millisecond)
			
			if err != nil {
				break
			}

			var transaction models.Transaction
			json.Unmarshal(msg.Value, &transaction)

			// Convert spending amount to respective point-type
			services.ConvertPoints(&transaction)

			// Retrieve and apply campaign bonus, if applicable
			services.ApplyApplicableCampaign(&transaction)

			fmt.Println("Transaction:", transaction)  // DEBUG
			transactions.Transactions = append(transactions.Transactions, transaction)
		}

		// If there are transactions, insert them into the DB and commit
		if len(transactions.Transactions) != 0 {
			_, err := collections.CreateTransactions(transactions)
			if err == nil {
				consumer.Commit()
			}
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