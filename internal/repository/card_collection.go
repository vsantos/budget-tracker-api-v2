package repository

import (
	"budget-tracker-api-v2/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// CardCollectionInterface defines a mongodb collection API to be posteriorly mocked
type CardCollectionInterface interface {
	CreateIndexes(ctx context.Context, indexes []string) error
	InsertOne(ctx context.Context, document interface{}) (id string, err error)
	FindOne(ctx context.Context, id string) (*model.Card, error)
	FindOneByFilter(ctx context.Context, filter bson.M) (*model.Card, error)
	DeleteOne(ctx context.Context, id string) (int64, error)
}
