package mongodb

import (
	"budget-tracker-api-v2/internal/model"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SpendCollectionConfig will implement mongodb collection functions
type SpendCollectionConfig struct {
	MongoCollection *mongo.Collection
}

// CreateIndexes will create mongodb indexes
func (c *SpendCollectionConfig) CreateIndexes(ctx context.Context, indexes []string) error {
	var indexModels []mongo.IndexModel
	for _, i := range indexes {
		indexModels = append(indexModels, mongo.IndexModel{
			Keys:    bson.D{{Key: i, Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	}

	_, err := c.MongoCollection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return err
	}

	return nil
}

// InsertOne will insert a document into mongodb
func (c *SpendCollectionConfig) InsertOne(ctx context.Context, document interface{}) (id string, err error) {
	if document == nil {
		return "", errors.New("empty document to register")
	}

	r, err := c.MongoCollection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", r.InsertedID), nil
}

// FindOne will find a Spend from collection
func (c *SpendCollectionConfig) FindOne(ctx context.Context, id string) (*model.Spend, error) {
	var emp model.Spend

	i, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("not able to get ID format")
	}

	filter := bson.M{"_id": i}
	err = c.MongoCollection.FindOne(ctx, filter).Decode(&emp)

	if err != nil {
		return nil, err
	}

	return &emp, nil
}

// DeleteOne will find a Spend from collection
func (c *SpendCollectionConfig) DeleteOne(ctx context.Context, id string) (int64, error) {
	i, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, errors.New("not able to get ID format")
	}

	filter := bson.M{"_id": i}
	result, err := c.MongoCollection.DeleteOne(ctx, filter)

	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
