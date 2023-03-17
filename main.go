package main

import (
	"fmt"
	"log"
	"time"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

/*
 * TODO:
 * 1. Migrate POST method from controller.transaction.go after Gabriel makes it ATOMIC
 * 2. Implement VALIDATION checks (TBC)
 * 3. Implement reading from Kafka
 */

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


func main() {

	// Setting up a connection with kafka
	// server := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        "localhost:9092",
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
	consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}

	run := true

	// Run a infinite loop that constantly checks for messages
	for run {

		msg, err := consumer.ReadMessage(time.Second)

		if err == nil {
			// TODO: Process transaction
			
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			// Only commit after successfully processed the message
			consumer.CommitMessage(msg)
		} else if err != nil {
			// TODO Handle error
		}
	}

	consumer.Close()

}
