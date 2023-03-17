package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DBinstance()

func DBinstance() (client *mongo.Client) {

	user := os.Getenv("MONGO_USERNAME")
	pass := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")

	conn := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)
	if port == "" || port == "443" {
		fmt.Println("Using mongo+srv config")
		conn = fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority", user, pass, host)
	}
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(conn).SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// connect to mongodb
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// initialise indexes
	InitIndexes(client)

	return client
}

func InitIndexes(client *mongo.Client) {

	// transactions_transactions_-1 index
	transactionCollection := OpenCollection(client, "transactions")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"transaction_id", -1}},
		Options: options.Index().SetUnique(true),
	}
	indexCreated, err := transactionCollection.Indexes().CreateOne(context.Background(), indexModel)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created Index %s\n", indexCreated)
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("loyalty").Collection(collectionName)

	return collection
}
