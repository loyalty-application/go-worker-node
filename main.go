package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/loyalty-application/go-worker-node/models"
	"github.com/loyalty-application/go-worker-node/config"
	"github.com/loyalty-application/go-worker-node/services"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var transactionCollection *mongo.Collection = config.OpenCollection(config.Client, "transactions")

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

		for i := 0; i < 50000; i++ {
			msg, err := consumer.ReadMessage(time.Millisecond)
			
			if err != nil { break }

			var transaction models.Transaction
			json.Unmarshal(msg.Value, &transaction)
			services.ConvertPoints(&transaction)
			fmt.Println(transaction)  // DEBUG
			transactions.Transactions = append(transactions.Transactions, transaction)
		}

		// If there are transactions, insert them into the DB and commit
		if len(transactions.Transactions) != 0 {
			_, err := CreateTransactions(transactions)
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


func CreateTransactions(transactions models.TransactionList) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// convert from slice of struct to slice of interface
	t := make([]interface{}, len(transactions.Transactions))
	for i, v := range transactions.Transactions {
		t[i] = v
	}

	// convert from slice of interface to mongo's bulkWrite model
	models := make([]mongo.WriteModel, 0)
	for _, doc := range t {
		models = append(models, mongo.NewInsertOneModel().SetDocument(doc))
	}
	
	// If an error occurs during the processing of one of the write operations, MongoDB
	// will continue to process remaining write operations in the list.
	bulkWriteOptions := options.BulkWrite().SetOrdered(false)
	// log.Println("Bulk Writing", models)
	result, err = transactionCollection.BulkWrite(ctx, models, bulkWriteOptions)
    if err != nil && !mongo.IsDuplicateKeyError(err) {
        log.Println(err.Error())
		return result, err
    }

	return result, err
}