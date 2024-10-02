package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main-service/config"
)

func Connect(ctx context.Context, c *config.Config) (db *mongo.Database, err error) {
	clientOptions := options.Client().ApplyURI(c.MongoDBDSN())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb: %w", err)
	}
	//Ping
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error pinging mongodb: %w", err)
	}
	return client.Database(c.MongoDB.DBName), nil
}

func Disconnect(ctx context.Context, client *mongo.Client) error {
	return client.Disconnect(ctx)
}
