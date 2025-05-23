package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"hw_14/internal/llogger"
)

const (
	collectionName = "kv"
	dbName         = "app"
	DbUri          = "mongodb://root:root@localhost:27017"
)

var logger = llogger.SetupLogger()

func createClient(ctx context.Context, opts *options.ClientOptions) *mongo.Client {
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(fmt.Errorf("mongodb connect failed: %v", err))
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(fmt.Errorf("mongodb ping failed: %v", err))
	}
	logger.Debug("mongodb connected")
	return client
}

func EnsureCollection(ctx context.Context, client *mongo.Client, dbName, collectionName string) error {
	db := client.Database(dbName)

	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}

	for _, name := range collections {
		if name == collectionName {
			logger.Debug("collection already exists")
			return nil
		}
	}

	err = db.CreateCollection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("failed to create collection: %v", err)
	}
	logger.Debug("collection created")
	return nil
}

func main() {
	ctx := context.Background()
	opts := options.Client()
	opts.ApplyURI(DbUri)
	client := createClient(ctx, opts)
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			panic(fmt.Errorf("mongodb disconnect failed: %v", err))
		}
	}(client, ctx)

	err := EnsureCollection(ctx, client, dbName, collectionName)
	if err != nil {
		panic(fmt.Errorf("failed to create collection: %v", err))
	}

}
