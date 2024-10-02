package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"storage-service/config"
)

func Connect(ctx context.Context, c *config.Config) (db *mongo.Database, err error) {
	clientOptions := options.Client().ApplyURI(c.MongoDBDSN())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	//Ping
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client.Database(c.MongoDB.DBName), nil
}
