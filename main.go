package main

import (
	"encoding/json"
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

	// Create a new Kafka Consumer
	consumer, err := getKafkaConsumer()
	defer consumer.Close()

	// Subscribe to kafka topic based on worker node type
	workerType := os.Getenv("WORKER_NODE_TYPE")
	topic := ""
	if workerType == "users" {
		log.Println("In here")
		topic = "users"
	} else if workerType == "transactions" {
		topic = "ftptransactions"
	}
	consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("consuming ... !! Type =", topic)
	if workerType == "users" {
		processUsers(consumer)
	} else if workerType == "transactions" {
		processTransactions(consumer)
	}
}

func processUsers(consumer *kafka.Consumer) {

	for {
		var userRecords models.UserRecordList
		var users models.UserList

		for i := 0; i < 20000; i++ {
			msg, err := consumer.ReadMessage(time.Second)
			
			if err != nil { break }

			var userRecord models.UserRecord
			json.Unmarshal(msg.Value, &userRecord)
			user, err := services.GetUserFromRecord(userRecord)

			// Debug
			log.Println("User Record", userRecord)
			log.Println("User", user)

			userRecords.UserRecords = append(userRecords.UserRecords, userRecord)
			users.Users = append(users.Users, user)
		}

		// If there are transactions, insert them into the DB and commit
		if len(userRecords.UserRecords) != 0 {
			// collections.CreateTransactions(transactions)
			log.Println("Committing", userRecords.UserRecords)
			consumer.Commit()
		}

	}
}

func processTransactions(consumer *kafka.Consumer) {

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