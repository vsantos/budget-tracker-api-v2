package repository

import (
	"budget-tracker-api-v2/internal/model"
	"context"
)

// UserCollectionInterface defines a mongodb collection API to be posteriorly mocked
type UserCollectionInterface interface {
	CreateIndexes(ctx context.Context, indexes []string) error
	InsertOne(ctx context.Context, document interface{}) (id string, err error)
	FindOne(ctx context.Context, id string) (*model.User, error)
	DeleteOne(ctx context.Context, id string) (int64, error)
}
