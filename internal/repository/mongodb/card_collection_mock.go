package mongodb

import (
	"budget-tracker-api-v2/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CardMockCollectionConfig will implement mongodb collection functions
type CardMockCollectionConfig struct {
	Error error
}

// CreateIndexes will create mongodb indexes
func (c *CardMockCollectionConfig) CreateIndexes(ctx context.Context, indexes []string) error {
	return nil
}

// InsertOne will insert a document into mongodb
func (c *CardMockCollectionConfig) InsertOne(ctx context.Context, document interface{}) (id string, err error) {
	if c.Error != nil {
		return "", c.Error
	}

	return "66f1cca3c37c733c4ada103d", nil
}

// FindOne will insert a document into mongodb
func (c *CardMockCollectionConfig) FindOne(ctx context.Context, id string) (*model.Card, error) {

	return &model.Card{
		ID: primitive.NewObjectID(),
		// Login: "mockcard",
		// Name:  "Mock Card Torres",
		// Email: "mock.card@gmail.com",
	}, nil
}

// FindOne will insert a document into mongodb
func (c *CardMockCollectionConfig) FindOneByFilter(ctx context.Context, filter bson.M) (*model.Card, error) {

	return &model.Card{
		ID: primitive.NewObjectID(),
	}, nil
}

// DeleteOne will insert a document into mongodb
func (c *CardMockCollectionConfig) DeleteOne(ctx context.Context, id string) (int64, error) {
	return 1, nil
}
