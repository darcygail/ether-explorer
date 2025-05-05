package store

import (
	"context"
	"strings"
	"time"

	"github.com/darcygail/ether-explorer/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveAsset(ctx context.Context, acc *schema.Account) error {

	collection := client.Database(DBName).Collection(AssetCollectionName)
	acc.UpdatedAt = time.Now().Unix()
	filter := bson.M{"address": strings.ToLower(acc.Address)}
	update := bson.M{
		"$set": bson.M{
			"balance":    acc.Balance,
			"assets":     acc.Assets,
			"updated_at": acc.UpdatedAt,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err

}

func GetAsset(ctx context.Context, address string) (schema.Asset, error) {
	collection := client.Database(DBName).Collection(AssetCollectionName)
	var account schema.Asset // 你定义的结构体
	err := collection.FindOne(ctx, bson.M{"address": strings.ToLower(address)}).Decode(&account)
	return account, err
}
