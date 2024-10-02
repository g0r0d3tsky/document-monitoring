package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type StorageMongo struct {
	collection *mongo.Collection
}

func NewStorageMongo(client *mongo.Client, dbName, collectionName string) *StorageMongo {
	collection := client.Database(dbName).Collection(collectionName)
	return &StorageMongo{collection: collection}
}

func (r *StorageMongo) SaveTextContent(ctx context.Context, fileContent []byte, filename string) error {
	content := map[string]interface{}{
		"filename":  filename,
		"text":      string(fileContent),
		"createdAt": time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, content)
	if err != nil {
		return fmt.Errorf("error inserting document into MongoDB: %w", err)
	}

	return nil
}
