package mongodb

import (
	"budget-tracker-api-v2/internal/model"
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

// FindOneBy will find a document based on field
func (c *UserMockCollectionConfig) FindOneBy(ctx context.Context, login string) (*model.User, error) {
	return &model.User{
		ID:        primitive.NewObjectID(),
		Login:     login,
		Firstname: "Mock User",
		Lastname:  "Torres",
		Email:     "mock.user@gmail.com",
		Password:  "$2a$10$HOrmuqyfwr575K4P9tjQXe0QKqbddMA/KFZ.YZhWVKPLMUF3LS4gi",
	}, nil
}

// FindOne will find a document based on ID
func (c *UserMockCollectionConfig) FindOne(ctx context.Context, id string) (*model.User, error) {

	return &model.User{
		ID:        primitive.NewObjectID(),
		Login:     "mockuser",
		Firstname: "Mock User",
		Lastname:  "Torres",
		Email:     "mock.user@gmail.com",
		// Salted password
		Password: "$2a$10$HOrmuqyfwr575K4P9tjQXe0QKqbddMA/KFZ.YZhWVKPLMUF3LS4gi",
	}, nil
}

// DeleteOne will insert a document into mongodb
func (c *UserMockCollectionConfig) DeleteOne(ctx context.Context, id string) (int64, error) {
	return 1, nil
}
