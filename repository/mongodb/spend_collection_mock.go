package mongodb

import (
	"budget-tracker-api-v2/model"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SpendMockCollectionConfig will implement mongodb collection functions
type SpendMockCollectionConfig struct {
	Error error
}

// CreateIndexes will create mongodb indexes
func (c *SpendMockCollectionConfig) CreateIndexes(ctx context.Context, indexes []string) error {
	return nil
}

// InsertOne will insert a document into mongodb
func (c *SpendMockCollectionConfig) InsertOne(ctx context.Context, document interface{}) (id string, err error) {
	if c.Error != nil {
		return "", c.Error
	}

	return "66f1cca3c37c733c4ada103d", nil
}

// FindOne will insert a document into mongodb
func (c *SpendMockCollectionConfig) FindOne(ctx context.Context, id string) (*model.Spend, error) {

	return &model.Spend{
		ID: primitive.NewObjectID(),
		// Login: "mockspend",
		// Name:  "Mock Spend Torres",
		// Email: "mock.spend@gmail.com",
	}, nil
}

// DeleteOne will insert a document into mongodb
func (c *SpendMockCollectionConfig) DeleteOne(ctx context.Context, id string) (int64, error) {
	return 1, nil
}
