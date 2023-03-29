package collections

import (
	"context"
	"log"
	"time"

	"github.com/loyalty-application/go-worker-node/config"
	"github.com/loyalty-application/go-worker-node/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cardCollection *mongo.Collection = config.OpenCollection(config.Client, "cards")

func RetrieveSpecificCard(cardId string) (result models.Card, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "card_id", Value: cardId}}
	err = cardCollection.FindOne(ctx, filter).Decode(&result)
	
	return result, err
}

func CreateCard(card models.Card) (result *mongo.InsertOneResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	result, err = cardCollection.InsertOne(ctx, card)

	return result, err
}

func CreateCards(cards models.CardList) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// convert from slice of struct to slice of interface
	t := make([]interface{}, len(cards.Cards))
	for i, v := range cards.Cards {
		t[i] = v
	}

	// convert from slice of interface to mongo's bulkWrite model
	models := make([]mongo.WriteModel, 0)
	for _, doc := range t {
		models = append(models, mongo.NewInsertOneModel().SetDocument(doc))
	}
	
	// If an error occurs during the processing of one of the write operations, MongoDB
	// will continue to process remaining write operations in the list.
	bulkWriteOptions := options.BulkWrite().SetOrdered(false)
	// log.Println("Bulk Writing", models)
	result, err = cardCollection.BulkWrite(ctx, models, bulkWriteOptions)
    if err != nil && !mongo.IsDuplicateKeyError(err) {
        log.Println(err.Error())
		panic(err)
    }

	return result, err

}