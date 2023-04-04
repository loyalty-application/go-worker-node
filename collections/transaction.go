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

var transactionCollection *mongo.Collection = config.OpenCollection(config.Client, "transactions")
var unprocessedCollection *mongo.Collection = config.OpenCollection(config.Client, "unprocessed")

func CreateTransactions(transactions models.TransactionList) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// convert from slice of struct to slice of interface
	t := make([]interface{}, len(transactions.Transactions))
	for i, v := range transactions.Transactions {
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
	result, err = transactionCollection.BulkWrite(ctx, models, bulkWriteOptions)
    if err != nil && !mongo.IsDuplicateKeyError(err) {
        log.Println(err.Error())
		return result, err
    }

	return result, err
}

func RetrieveCardValuesFromTransaction(cardId string) (result float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{"card_id": cardId}},
		{"$group": bson.M{
			"_id": nil,
			"totalPoints": bson.M{"$sum": "$points"},
			"totalMiles": bson.M{"$sum": "$miles"},
			"totalCashback": bson.M{"$sum": "$cash_back"},
		}},
	}
	
	cursor, err := transactionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err.Error())
		return result, err
	}
	
	var temp struct {
		TotalPoints   float64 `bson:"totalPoints"`
		TotalMiles    float64 `bson:"totalMiles"`
		TotalCashback float64 `bson:"totalCashback"`
	}

	if cursor.Next(context.Background()) {
        
		if err = cursor.Decode(&temp); err != nil {
			log.Println(err.Error())
			return result, err
		}
	}

	result += temp.TotalCashback + temp.TotalMiles + temp.TotalPoints

	return result, err
}

func DeleteUnprocessedByTransactionId(transactionIdList []string) (result *mongo.DeleteResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"transaction_id": bson.M{"$in": transactionIdList}}

	result, err = unprocessedCollection.DeleteMany(ctx, filter)

	return result, err
}