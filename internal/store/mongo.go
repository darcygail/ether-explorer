package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const DBName = "ether_explorer"
const AssetCollectionName = "account_assets"

func InitDB(mongoUri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		return err
	}
	// create address index
	collection := client.Database(DBName).Collection(AssetCollectionName)
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]int{"address": 1},
		Options: options.Index().SetUnique(true),
	})

	client = c
	return nil
}
