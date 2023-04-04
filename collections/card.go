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

func CreateCards(cards []models.Card) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// convert from slice of interface to mongo's bulkWrite model
	models := make([]mongo.WriteModel, 0)
	for _, card := range cards {
		update := bson.D{{"$set", bson.D{{"card_pan", card.CardPan},
										{"card_type", card.CardType},
										{"value_type", card.ValueType},
										{"short_card_pan,", card.ShortCardPan},
										{"user_id", card.UserId}}},
						 {"$max", bson.D{{"value", card.Value}}}}

		upsert := mongo.NewUpdateOneModel().SetFilter(bson.M{"card_id":card.CardId}).SetUpdate(update).SetUpsert(true)

		models = append(models, upsert)
	}
	
	// If an error occurs during the processing of one of the write operations, MongoDB
	// will continue to process remaining write operations in the list.
	bulkWriteOptions := options.BulkWrite().SetOrdered(false)
	result, err = cardCollection.BulkWrite(ctx, models, bulkWriteOptions)
    if err != nil {
        log.Println(err.Error())
    }

	return result, err

}

func UpdateCardValues(cardMap map[string]float64) (result *mongo.BulkWriteResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	models := make([]mongo.WriteModel, 0)
	for cardId, value := range cardMap {
		update := bson.D{{"$inc", bson.D{{"value", value}}}}

		upsert := mongo.NewUpdateOneModel().SetFilter(bson.M{"card_id":cardId}).SetUpdate(update).SetUpsert(true)

		models = append(models, upsert)
	}

	// Create a new bulk write options instance
	opts := options.BulkWrite().SetOrdered(false)

	// Execute the bulk write operation with the writes array and options
	result, err = cardCollection.BulkWrite(ctx, models, opts)
	if err != nil {
		log.Println("Error =", err.Error())
	}

	return result, err
}

func RetrieveEmailFromCard(cardId string) (email string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	pipeline := []bson.M{
        bson.M{
            "$match": bson.M{"card_id": cardId},
        },
        bson.M{
            "$lookup": bson.M{
                "from":         userCollection.Name(),
                "localField":   "user_id",
                "foreignField": "user_id",
                "as":           "user",
            },
        },
        bson.M{
            "$project": bson.M{
                "_id":    0,
                "email": bson.M{"$arrayElemAt": []interface{}{"$user.email", 0}},
            },
        },
    }

    cursor, err := cardCollection.Aggregate(ctx, pipeline)
    if err != nil {
        return "", err
    }

    var result struct {
        Email string `bson:"email"`
    }

    if cursor.Next(ctx) {
        err = cursor.Decode(&result)
        if err != nil {
            return "", err
        }
    } else {
        return "", mongo.ErrNoDocuments
    }

	return result.Email, err
}