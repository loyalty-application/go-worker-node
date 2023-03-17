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
	// err = client.Ping(context.Background(), nil)
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
	user := os.Getenv("MONGO_USERNAME")
	pass := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")

	mongoURL := "mongodb://" + user + ":" + pass + "@" + host + ":" + port
	// get Mongo db Collection using environment variables.
	dbName := "loyalty"
	collectionName := "transactions"
	collection := getMongoCollection(mongoURL, dbName, collectionName)
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

	defer consumer.Close()


	topic := "ftptransactions"

	consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatal(err)
	}



	fmt.Println("start consuming ... !!")
	for {

		msg, err := consumer.ReadMessage(time.Second)

		if err == nil {
			// TODO: Process transaction
			var transaction models.Transaction
			json.Unmarshal(msg.Value, &transaction)
			insertResult, err := collection.InsertOne(context.Background(), transaction)
			
			fmt.Println("Inserted a single document: ", insertResult.InsertedID)

			// Only commit after successfully processed the message
			if err == nil {

				consumer.CommitMessage(msg)
			}
		} else if err != nil {
			// TODO Handle error
		}
	}


}
