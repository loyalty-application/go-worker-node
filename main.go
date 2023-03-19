package main

import (
	"strconv"
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

	for {
		var transactions models.TransactionList

		for i := 0; i < 50000; i++ {
			msg, err := consumer.ReadMessage(time.Millisecond)
			
			if err != nil { break }

			var transaction models.Transaction
			json.Unmarshal(msg.Value, &transaction)
			fmt.Println("Begin converting points")
			convertPoints(&transaction)
			fmt.Println("Finished converting points")
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

// QUESTIONS:
// 1. How to determine online shopping?
// 2. How to determine exchange rate?

func getExchangeRate(currency string) float64 {
	// TODO: Implement an exchange rate getter
	return 1.0
}

func in(arr []int, target int) bool {
	for _, item := range arr {
		if item == target { return true }
	}
	return false
}

func convertPoints(transaction *models.Transaction) {

	// MCC code categorization (SUBJECT TO CHANGES)
	var ONLINE_SHOPPING = []int{5999, 5964, 5691, 5311, 5411, 5399}
	var HOTEL = []int{7011}
	var EXCLUDED_MCCS = []int{6051, 9399, 6540}

	// If MCC is to be excluded, do not convert to points
	mcc, _ := strconv.Atoi(transaction.MCC)
	if in(EXCLUDED_MCCS, mcc) {
		return
	}

	// Get amount spent
	amountSpent := transaction.Amount * getExchangeRate(transaction.Currency)

	// Conversion: $$ -> POINTS
	if transaction.CardType == "scis_shopping" {
		pointsConversionRate := 4
		if in(ONLINE_SHOPPING, mcc) {  // Bonus for online shopping
			pointsConversionRate = 10
		}
		transaction.Points = amountSpent * float64(pointsConversionRate + 1)

	// Conversion: $$ -> CASHBACK
	} else if transaction.CardType == "scis_freedom" {
		cashBack := 0.005 * amountSpent
		if amountSpent > 500 {
			cashBack += 0.01 * amountSpent
		}
		if amountSpent > 2000 {
			cashBack += 0.03 * amountSpent
		}
		transaction.CashBack = cashBack

	// Conversion: $$ -> MILES
	} else {
		transaction.Miles = 0
	}
}