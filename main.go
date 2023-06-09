package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/loyalty-application/go-worker-node/collections"
	"github.com/loyalty-application/go-worker-node/config"
	"github.com/loyalty-application/go-worker-node/models"
	"github.com/loyalty-application/go-worker-node/services"
	"github.com/loyalty-application/go-worker-node/testing"
)

func main() {

	// Connect to DB
	config.DBinstance()

	// TESTING (INSERT CAMPAIGNS)
	testing.AddCampaignsTest()
	log.Println("TEST: ADDED CAMPAIGNS")

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

	log.Println("im consuming ... !! Type =", topic)
	if workerType == "users" {
		processUsers(consumer)
	} else if workerType == "transactions" {
		processFtpTransactions(consumer)
	} else if workerType == "resttransactions" {
		processRestTransactions(consumer)
	}
}

func processUsers(consumer *kafka.Consumer) {

	for {
		var users models.UserList
		var cards models.CardList

		for i := 0; i < 20000; i++ {
			msg, err := consumer.ReadMessage(time.Second)

			if err != nil {
				break
			}

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

		}

		// If there are users / cards, insert them into the DB and commit
		// TODO Implement Goroutines here
		if len(users.Users) != 0 {
			log.Println("Appending Users, Len =", len(users.Users))
			collections.CreateUsers(users)
		}

		if len(cards.Cards) != 0 {
			log.Println("Appending Cards, Len =", len(cards.Cards))
			collections.CreateCards(cards.Cards)
		}

		if len(cards.Cards) != 0 || len(users.Users) != 0 {
			consumer.Commit()
		}
	}
}

func processFtpTransactions(consumer *kafka.Consumer) {

	const BATCH_SIZE int = 20000

	for {
		var transactions models.TransactionList

		// key = cardId, value = points / miles / cashback
		cardMap := make(map[string]float64)

		// Retrieve all campaigns
		allCampaigns, _ := collections.RetrieveAllCampaigns()
		notificationList := make([]models.Notification, 0)

		// Process transactions in batches
		for i := 0; i < BATCH_SIZE; i++ {

			msg, err := consumer.ReadMessage(time.Second)

			if err != nil {
				break
			}

			var transaction models.Transaction
			json.Unmarshal(msg.Value, &transaction)

			// Only apply points conversion for valid transaction
			if services.IsValidTransaction(&transaction) {
				services.ConvertPoints(&transaction)
				
				// Apply applicable campaigns
				campaign, hasCampaign := services.ApplyApplicableCampaign(&transaction, allCampaigns)

				// Create email notification
				if hasCampaign {
					message := "Hi, you have successfully qualified for a Campaign by " + campaign.Merchant + ". Campaign's description" + campaign.Description
					notificationList = append(notificationList, models.Notification{ CardId: transaction.CardId,
																					Message: message,})
				}
			}

			// Update CardMap
			cardMap[transaction.CardId] += transaction.Points + transaction.Miles + transaction.CashBack

			// Hash CardPan
			transaction.CardPan = services.HashCardPan(transaction.CardPan)

			// Add transaction into regardless of validity
			transactions.Transactions = append(transactions.Transactions, transaction)
		}

		// If there are transactions, insert them into the DB and commit
		if len(transactions.Transactions) != 0 {
			// Commit transaction
			collections.CreateTransactions(transactions)

			// Update card points after committing transactions (Upsert if necessary)
			// TODO Implement Goroutines here
			collections.UpdateCardValues(cardMap)

			// Send email notification, if any
			// log.Println(notificationList)
			// services.SendNotification(notificationList)

			// Testing
			log.Println("Processed", len(transactions.Transactions), " transactions")
			consumer.Commit()
		}
	}
}

func processRestTransactions(consumer *kafka.Consumer) {

	for {
		var transactions models.TransactionList

		// key = cardId, value = points / miles / cashback
		cardMap := make(map[string]float64)

		// Retrieve All Campaigns
		allCampaigns, _ := collections.RetrieveAllCampaigns()

		notificationList := make([]models.Notification, 0)

		transactionIdList := make([]string, 0)

		for i := 0; i < 20000; i++ {
			msg, err := consumer.ReadMessage(time.Second)

			if err != nil {
				break
			}

			var transaction models.Transaction
			json.Unmarshal(msg.Value, &transaction)


			// Only apply points conversion for valid transaction
			if services.IsValidTransaction(&transaction) {
				services.ConvertPoints(&transaction)
				
				// Apply applicable campaigns
				campaign, hasCampaign := services.ApplyApplicableCampaign(&transaction, allCampaigns)

				// Create email notification
				if hasCampaign {
					message := "Hi, you have successfully qualified for a Campaign by " + campaign.Merchant + ". Campaign's description" + campaign.Description
					notificationList = append(notificationList, models.Notification{ CardId: transaction.CardId,
																					Message: message,})
				}
			}

			// Update CardMap
			cardMap[transaction.CardId] += transaction.Points + transaction.Miles + transaction.CashBack

			// Hash CardPan
			transaction.CardPan = services.HashCardPan(transaction.CardPan)

			// Add transaction into regardless of validity
			transactions.Transactions = append(transactions.Transactions, transaction)

			// Add transaction into list
			transactionIdList = append(transactionIdList, transaction.TransactionId)
		}

		// If there are transactions, insert them into the DB and commit
		if len(transactions.Transactions) != 0 {
			// Commit transaction
			collections.CreateTransactions(transactions)

			// Delete from unprocessed collection
			result, err := collections.DeleteUnprocessedByTransactionId(transactionIdList)
			if err != nil {
				log.Println("Delete from Unprocessed Error:", err.Error())
			} else {
				log.Println("Deleted:", result.DeletedCount)
			}

			// Update card points after committing transactions (Upsert if necessary)
			// TODO Implement Goroutines here
			collections.UpdateCardValues(cardMap)

			// Send email notification, if any
			log.Println(notificationList)
			services.SendNotification(notificationList)

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
