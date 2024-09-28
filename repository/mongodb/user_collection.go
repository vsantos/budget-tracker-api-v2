package mongodb

import (
	"budget-tracker-api-v2/model"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserCollectionConfig will implement mongodb collection functions
type UserCollectionConfig struct {
	MongoCollection *mongo.Collection
}

// CreateIndexes will create mongodb indexes
func (c *UserCollectionConfig) CreateIndexes(ctx context.Context, indexes []string) error {
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
func (c *UserCollectionConfig) InsertOne(ctx context.Context, document interface{}) (id string, err error) {
	r, err := c.MongoCollection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", r.InsertedID), nil
}

// FindOne will find a User from collection
func (c *UserCollectionConfig) FindOne(ctx context.Context, id string) (*model.User, error) {
	var emp model.User

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