package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DBinstance()

func DBinstance() (client *mongo.Client) {
	godotenv.Load()

	user := os.Getenv("MONGO_USERNAME")
	pass := os.Getenv("MONGO_PASSWORD")
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")

	conn := fmt.Sprintf("mongodb://%s:%s@%s:%s", user, pass, host, port)

	replicaSetQueryString := "/?replicaSet=replica-set"
	tlsQueryString := ""
	secondaryQueryString := ""
	//var tlsConfig *tls.Config

	if os.Getenv("GIN_MODE") == "release" {
		replicaSetQueryString = "/?replicaSet=rs0"
		//tlsQueryString = "&tls=true"
		secondaryQueryString = "&readPreference=secondaryPreferred&retryWrites=false"

		//// configure tls
		//var filename = "rds-combined-ca-bundle.pem"
		//tlsConfig := new(tls.Config)
		//certs, err := ioutil.ReadFile(filename)

		//if err != nil {
		//fmt.Println("Failed to read CA file")
		//return
		//}

		//tlsConfig.RootCAs = x509.NewCertPool()
		//ok := tlsConfig.RootCAs.AppendCertsFromPEM(certs)

		//if !ok {
		//fmt.Println("Failed to append CA file")
		//return
		//}

		//if tlsConfig != nil {
		//fmt.Println("Successfully set TLS config")
		//clientOptions.SetTLSConfig(tlsConfig)
		//}

	}
	conn = fmt.Sprintf("%s%s%s%s", conn, replicaSetQueryString, tlsQueryString, secondaryQueryString)

	fmt.Printf("Attempting connection with: %s\n", conn)
	//serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	//clientOptions := options.Client().ApplyURI(conn).SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to mongodb
	//client, err := mongo.Connect(ctx, clientOptions)
	//if err != nil {
	//log.Fatal(err)
	//}
	client, err := mongo.NewClient(options.Client().ApplyURI(conn))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("Connecting to MongoDB ...")
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to cluster: %v", err)
	}

	fmt.Println("Success!")
	// Force a connection to verify our connection string

	fmt.Println("Pinging server ...")
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping cluster: %v", err)
	}
	fmt.Println("Success!")

	fmt.Println("Initialising indexes ...")
	// initialise indexes
	InitIndexes(client)
	fmt.Println("Success!")
	return client
}

func InitIndexes(client *mongo.Client) {

	// transactions_transactions_-1 index
	transactionCollection := OpenCollection(client, "transactions")

	transactionIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_id", Value: -1}},
		Options: options.Index().SetUnique(true),
	}
	transactionIndexCreated, err := transactionCollection.Indexes().CreateOne(context.Background(), transactionIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	// unprocessed_unprocessed-1 index
	unprocessedCollection := OpenCollection(client, "unprocessed")

	unprocessedIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "transaction_id", Value: -1}},
		Options: options.Index().SetUnique(true),
	}
	unprocessedIndexCreated, err := unprocessedCollection.Indexes().CreateOne(context.Background(), unprocessedIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	// campaigns_campaigns_-1 index
	campaignCollection := OpenCollection(client, "campaigns")

	campaignIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "campaign_id", Value: -1}},
		Options: options.Index().SetUnique(true),
	}

	campaignIndexCreated, err := campaignCollection.Indexes().CreateOne(context.Background(), campaignIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	// cards_cards_-1 index
	cardCollection := OpenCollection(client, "cards")

	cardIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "card_id", Value: -1}},
		Options: options.Index().SetUnique(true),
	}

	cardIndexCreated, err := cardCollection.Indexes().CreateOne(context.Background(), cardIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	// user_users_-1 index
	userCollection := OpenCollection(client, "users")

	userIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: -1}},
		Options: options.Index().SetUnique(true),
	}

	userIndexCreated, err := userCollection.Indexes().CreateOne(context.Background(), userIndexModel)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created Transaction Index %s\n", transactionIndexCreated)
	fmt.Printf("Created Unprocessed Index %s\n", unprocessedIndexCreated)
	fmt.Printf("Created Campaign Index %s\n", campaignIndexCreated)
	fmt.Printf("Created Card Index %s\n", cardIndexCreated)
	fmt.Printf("Created User Index %s\n", userIndexCreated)
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("loyalty").Collection(collectionName)

	return collection
}
