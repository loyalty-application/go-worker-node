package collections

import (
	"context"
	"time"

	"github.com/loyalty-application/go-worker-node/config"
	"github.com/loyalty-application/go-worker-node/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = config.OpenCollection(config.Client, "users")

func RetrieveSpecificUser(email string) (result models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}

	err = userCollection.FindOne(ctx, filter).Decode(&result)
	
	return result, err
}

func CreateUser(user models.User) (result *mongo.UpdateResult, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "user_id", Value: user.UserID}}
	update := bson.D{{Key: "$set", Value: user}}
	opts := options.Update().SetUpsert(true)

	result, err = userCollection.UpdateOne(ctx, filter, update, opts)

	return result, err

}