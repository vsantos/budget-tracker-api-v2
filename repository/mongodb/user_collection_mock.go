package mongodb

import (
	"budget-tracker-api-v2/model"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserMockCollectionConfig will implement mongodb collection functions
type UserMockCollectionConfig struct {
	Error error
}

// CreateIndexes will create mongodb indexes
func (c *UserMockCollectionConfig) CreateIndexes(ctx context.Context, indexes []string) error {
	return nil
}

// InsertOne will insert a document into mongodb
func (c *UserMockCollectionConfig) InsertOne(ctx context.Context, document interface{}) (id string, err error) {
	if c.Error != nil {
		return "", c.Error
	}

	return "66f1cca3c37c733c4ada103d", nil
}

// FindOne will insert a document into mongodb
func (c *UserMockCollectionConfig) FindOne(ctx context.Context, id string) (*model.User, error) {

	return &model.User{
		ID:    primitive.NewObjectID(),
		Login: "mockuser",
		Name:  "Mock User Torres",
		Email: "mock.user@gmail.com",
	}, nil
}

// DeleteOne will insert a document into mongodb
func (c *UserMockCollectionConfig) DeleteOne(ctx context.Context, id string) (int64, error) {
	return 1, nil
}
