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

var userCollection *mongo.Collection = config.OpenCollection(config.Client, "users")

func RetrieveSpecificUser(email string) (result models.User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}

	err = userCollection.FindOne(ctx, filter).Decode(&result)

	return result, err
}

func CreateUser(user models.User) (result *mongo.UpdateResult, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	filter := bson.D{{Key: "user_id", Value: user.UserID}}
	update := bson.D{{Key: "$set", Value: user}}
	opts := options.Update().SetUpsert(true)

	result, err = userCollection.UpdateOne(ctx, filter, update, opts)

	return result, err

}

func CreateUsers(users models.UserList) (result interface{}, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// convert from slice of struct to slice of interface
	t := make([]interface{}, len(users.Users))
	for i, v := range users.Users {
		userType := "USER"
		v.UserType = &userType
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
	result, err = userCollection.BulkWrite(ctx, models, bulkWriteOptions)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		log.Println(err.Error())
		panic(err)
	}

	return result, err

}
