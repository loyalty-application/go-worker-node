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
	count := 0  // counter to check messages consumed

	for {
		var transactions models.TransactionList

		for i := 0; i < 50000; i++ {
			msg, err := consumer.ReadMessage(time.Millisecond)
			
			if err != nil {
				break
			}

			var transaction models.Transaction
			json.Unmarshal(msg.Value, &transaction)
			fmt.Println("Card Type:", transaction.CardType)
			transactions.Transactions = append(transactions.Transactions, transaction)
			count += 1
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
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// convert from slice of struct to slice of interface
	t := make([]interface{}, len(transactions.Transactions))
	for i, v := range transactions.Transactions {
		t[i] = v
	}

	// Setting write permissions
	wc := writeconcern.New(writeconcern.WMajority())
	txnOpts := options.Transaction().SetWriteConcern(wc)

	// Start new session
	session, err := config.Client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())

	// Start transaction
	if err = session.StartTransaction(txnOpts); err != nil {
		return nil, err
	}
	log.Println("Transaction Start without errors")

	// Insert documents in the current session
	log.Println("Before Insert")
	result, err = transactionCollection.InsertMany(mongo.NewSessionContext(context.Background(),session), t)
	log.Println("After Insert")
	defer cancel()

	if err != nil {
		log.Println("Insert Many Error = ", err.Error())
		// Abort session if got error
		session.AbortTransaction(context.Background())
		// log.Println("Aborted Transaction")
		return result, err
	}

	// Commit documents if no error
	err = session.CommitTransaction(context.Background())

	return result, err
}


SCIS Shopping -- 
func convertPoints(transaction model.Transaction) {

	// MCC code categorization (SUBJECT TO CHANGES)
	const ONLINE_SHOPPING := []int{5999, 5964, 5691, 5311, 5411, 5399}
	const HOTEL := []int{7011}

	// Identify foreign or not foreign spend type

	// Get amount spent

	// Get card type

	// Convert according to card-point-type
}