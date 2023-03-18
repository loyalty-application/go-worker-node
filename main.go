package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
	"github.com/loyalty-application/go-worker-node/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoCollection(mongoURL, dbName, collectionName string) *mongo.Collection {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB ... !!")

	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	return collection
}

func main() {
	godotenv.Load(".env")
	// user := os.Getenv("MONGO_USERNAME")
	// pass := os.Getenv("MONGO_PASSWORD")
	// host := os.Getenv("MONGO_HOST")
	// port := os.Getenv("MONGO_PORT")

	// mongoURL := "mongodb://" + user + ":" + pass + "@" + host + ":" + port + "/?replicaSet=replica-set"
	// // get Mongo db Collection using environment variables.
	// dbName := "loyalty"
	// collectionName := "transactions"
	// collection := getMongoCollection(mongoURL, dbName, collectionName)
	server := os.Getenv("KAFKA_BOOTSTRAP_SERVER")
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        server,
		"group.id":                 "FtpWorkerGroup",
		"client.id":                "FtpProcessing",
		"enable.auto.commit":       true,
		"enable.auto.offset.store": true,
		"auto.offset.reset":        "earliest",
		"isolation.level":          "read_committed",
	})

	defer consumer.Close()


	topic := "ftptransactions"

	consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}



	fmt.Println("start consuming ... !!")
	// counter to check messages consumed
	count := 0
	for {

		var transactions models.TransactionList
		// var msglast *kafka.Message

		for i := 0; i < 30000; i++ {
			msg, err := consumer.ReadMessage(time.Millisecond)
			

			if err == nil {
				var transaction models.Transaction
				json.Unmarshal(msg.Value, &transaction)
				transactions.Transactions = append(transactions.Transactions, transaction)
				count += 1
			} else {
				fmt.Println(count)
			}
			
		}
		fmt.Println(count)
	}

}
