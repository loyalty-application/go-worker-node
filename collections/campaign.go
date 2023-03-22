package collections

import (
	"context"
	"time"

	"github.com/loyalty-application/go-worker-node/models"
	"github.com/loyalty-application/go-worker-node/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var campaignCollection *mongo.Collection = config.OpenCollection(config.Client, "campaigns")


// TODO: Replace with a RetrieveActiveCampaigns
func RetrieveAllCampaigns() (campaigns []models.Campaign, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Find().SetSort(bson.M{"start_date": 1})
	cursor, err := campaignCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, &campaigns)

	return campaigns, err
}