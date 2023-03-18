package main

import (
	"os"
	"fmt"
	"log"
	"time"
	// "net/http"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	// "go.mongodb.org/mongo-driver/mongo"
	// "github.com/gin-gonic/gin"
	// "github.com/loyalty-application/go-worker-node/models"
	// "github.com/loyalty-application/go-worker-node/config"
)

/*
 * TODO:
 * 1. Migrate POST method from controller.transaction.go after Gabriel makes it ATOMIC
 * 2. Implement VALIDATION checks (TBC)
 * 3. Implement reading from Kafka
 */

// var transactionCollection *mongo.Collection = config.OpenCollection(config.Client, "transactions")

// func CreateTransactions(userId string, transactions models.TransactionList) (result *mongo.InsertManyResult, err error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// convert from slice of struct to slice of interface
// 	t := make([]interface{}, len(transactions.Transactions))
// 	for i, v := range transactions.Transactions {
// 		v.UserId = userId
// 		t[i] = v
// 	}

// 	result, err = transactionCollection.InsertMany(ctx, t)
// 	return result, err
// }

// type TransactionController struct{}

// @Summary Create Transactions for User
// @Description Create transaction records
// @Tags    transaction
// @Accept  application/json
// @Produce application/json
// @Param   Authorization header string true "Bearer eyJhb..."
// @Param   user_id path string true "user's id"
// @Param   request body models.TransactionList true "transactions"
// @Success 200 {object} []models.Transaction
// @Failure 400 {object} models.HTTPError
// @Router  /transaction/{user_id} [post]
// func (t TransactionController) PostTransactions(c *gin.Context) {
// 	userId := c.Param("userId")
// 	if userId == "" {
// 		c.JSON(http.StatusBadRequest, models.HTTPError{http.StatusBadRequest, "Invalid User Id"})
// 		return
// 	}

// 	data := new(models.TransactionList)
// 	err := c.BindJSON(data)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, models.HTTPError{http.StatusBadRequest, "Invalid Transaction Object" + err.Error()})
// 		return
// 	}

// 	// TODO: make this operation atomic https://www.mongodb.com/docs/drivers/go/current/fundamentals/transactions/
// 	result, err := CreateTransactions(userId, *data)
// 	if err != nil {
// 		msg := "Invalid Transactions"
// 		if mongo.IsDuplicateKeyError(err) {
// 			msg = "transaction_id already exists"
// 		}
// 		c.JSON(http.StatusBadRequest, models.HTTPError{http.StatusBadRequest, msg})
// 		return
// 	}

// 	c.JSON(http.StatusOK, result)
// }

const DEBOUNCE_COUNTER = 6
const COMMIT_INTERVAL = 5

func main() {

	// Setting up a connection with kafka
	server := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        server,
		"group.id":                 "FtpWorkerGroup",
		"client.id":                "FtpProcessing",
		"enable.auto.commit":       false,
		"enable.auto.offset.store": false,
		"auto.offset.reset":        "earliest",
		"isolation.level":          "read_committed",
	})

	// Creating a topic categorization
	topic := "ftptransactions"

	// Subscribe to the message broker with decided topic
	log.Println("Subscribing")
	err = consumer.Subscribe(topic, nil)
	log.Println("Past Subscribe")
	if err != nil {
		log.Fatal(err)
	}

	timer := DEBOUNCE_COUNTER
	messageConsumed := 0
	var msg *kafka.Message = nil
	var prevMsg *kafka.Message = nil

	// Run a infinite loop that constantly checks for messages
	for true {
		// fmt.Println("Hello World")
		log.Println("========", timer)
		timer--
		prevMsg = msg
		msg, err := consumer.ReadMessage(time.Second)
		// fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

		if err == nil {
			// TODO: Process transaction
			
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			timer = DEBOUNCE_COUNTER  // Reset the timer
			messageConsumed++

			// Only commit after successfully processed the message
			if messageConsumed == COMMIT_INTERVAL {
				consumer.CommitMessage(msg)
				fmt.Println("Committing", messageConsumed, "messages")
				messageConsumed = 0
			}
		
		// No message to consume, commit once timer timeouts
		} else if err != nil {
			if timer == 0 {
				if prevMsg != nil {
					consumer.CommitMessage(prevMsg)
				}
				fmt.Println("DEB TIMEOUT: Committing", messageConsumed, "messages")
				timer = DEBOUNCE_COUNTER
				messageConsumed = 0
			}
		}
	}

	consumer.Close()

}
