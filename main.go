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
	log.Println(workerType)
	topic := ""
	if workerType == "users" {
		log.Println("In here")
		topic = "users"
	} else if workerType == "transactions" {
		topic = "ftptransactions"
	} else if workerType == "resttransactions" {
		topic = "resttransactions"
	}
	consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("world consuming ... !! Type =", topic)
	if workerType == "users" {
		processUsers(consumer)
	} else if workerType == "transactions" || workerType == "resttransactions" {
		processTransactions(consumer)
	}
}

func processUsers(consumer *kafka.Consumer) {

	for {
		var users models.UserList
		var cards models.CardList

		for i := 0; i < 20000; i++ {
			msg, err := consumer.ReadMessage(time.Second)
			
			if err != nil { break }

			var userRecord models.UserRecord
			json.Unmarshal(msg.Value, &userRecord)
			user, err := services.GetUserFromRecord(userRecord)

			if err == nil {
				users.Users = append(users.Users, user)
			}
			
			card, err := services.GetCardFromRecord(userRecord)
			if err == nil {
				// log.Println("Card =", card)
				cards.Cards = append(cards.Cards, card)
			}

			// // Debug
			// log.Println("User Record", userRecord)
			// log.Println("User", user)
			// log.Println("Card", card)
			// log.Println("Loop", i)
		}
		
		// If there are users / cards, insert them into the DB and commit
		if len(users.Users) != 0 {
			log.Println("Appending Users, Len =", len(users.Users))
			collections.CreateUsers(users)
		}

		if len(cards.Cards) != 0 {
			log.Println("Appending Cards, Len =", len(cards.Cards))
			collections.CreateCards(cards)
		}

		if len(cards.Cards) != 0 || len(users.Users) != 0 {
			consumer.Commit()
		}
	}
}

func processTransactions(consumer *kafka.Consumer) {

	for {
		var transactions models.TransactionList
		cardSet := map[string]struct{}{}

		for i := 0; i < 20000; i++ {
			msg, err := consumer.ReadMessage(time.Second)
			
			if err != nil {
				break
			}

			var transaction models.Transaction
			json.Unmarshal(msg.Value, &transaction)

			// Convert spending amount to respective point-type
			services.ConvertPoints(&transaction)
			
			// Add cardId used to set
			cardSet[transaction.CardId] = struct{}{}

			transactions.Transactions = append(transactions.Transactions, transaction)
		}

		// If there are transactions, insert them into the DB and commit
		if len(transactions.Transactions) != 0 {
			// Commit transaction
			collections.CreateTransactions(transactions)

			// Convert set of cards to slice of cards
			cardIdList := make([]string, 0)
			for cardId, _ := range cardSet {
				cardIdList = append(cardIdList, cardId)
			}
			log.Println("Card Id List =", cardIdList)

			// Update card points after committing transactions (Upsert if necessary)
			services.UpdateCardValues(cardIdList)

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