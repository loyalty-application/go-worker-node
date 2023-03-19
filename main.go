package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	// "net/http"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/loyalty-application/go-worker-node/models"
	"github.com/loyalty-application/go-worker-node/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

var transactionCollection *mongo.Collection = config.OpenCollection(config.Client, "transactions")

func main() {

	config.DBinstance()
	server := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        server,
		"group.id":                 "FtpWorkerGroup",
		"client.id":                "FtpProcessing",
		"enable.auto.commit":       false,
		"enable.auto.offset.store": true,
		"auto.offset.reset":        "earliest",
		"isolation.level":          "read_committed",
	})

	defer consumer.Close()


	topic := "ftptransactions"

	// Subscribe to the message broker with decided topic
	log.Println("Subscribing")
	err = consumer.Subscribe(topic, nil)
	log.Println("Past Subscribe")
	if err != nil {
		log.Fatal(err)
	}



	fmt.Println("start consuming ... !!")
	// counter to check messages consumed
	count := 0
	for {

		var transactions models.TransactionList

		for i := 0; i < 20000; i++ {
			msg, err := consumer.ReadMessage(time.Millisecond)
			

			if err == nil {
				var transaction models.Transaction
				json.Unmarshal(msg.Value, &transaction)
				transactions.Transactions = append(transactions.Transactions, transaction)
				count += 1
			} else {
				break
			}
			
		}

		if len(transactions.Transactions) != 0 {
			_, err := CreateTransactions(transactions)
			if err == nil {
				consumer.Commit()
			} else {
				log.Println(err.Error())
			}
		}
	}

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
    if err != nil {
        log.Println(err.Error())
		return result, err
    }

	return result, err
}