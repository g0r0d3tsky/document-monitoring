package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"main-service/config"
	"main-service/internal/usecase/model"
	"time"
)

type StorageMongo struct {
	collection *mongo.Collection
}

func NewStorageMongo(client *mongo.Client, cfg config.Config) *StorageMongo {
	collection := client.Database(cfg.MongoDB.DBName).Collection(cfg.MongoDB.CollectionName)
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

func (r *StorageMongo) GetContentByFilename(ctx context.Context, filename string) (*model.Content, error) {
	var content model.Content

	filter := bson.M{"filename": filename}
	err := r.collection.FindOne(ctx, filter).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no content found with filename: %s", filename)
		}
		return nil, fmt.Errorf("error finding document by filename: %w", err)
	}

	return &content, nil
}

func (r *StorageMongo) DeleteContentByFilename(ctx context.Context, filename string) error {
	filter := bson.M{"filename": filename}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting document from MongoDB: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no content found with filename: %s", filename)
	}

	return nil
}
