package repository

import (
	"budget-tracker-api-v2/internal/model"
	"context"
)

// TransactionCollectionInterface defines a mongodb collection API to be posteriorly mocked
type TransactionCollectionInterface interface {
	CreateIndexes(ctx context.Context, indexes []string) error
	InsertOne(ctx context.Context, t *model.Transaction) (transaction *model.Transaction, err error)
	FindOne(ctx context.Context, id string) (*model.Transaction, error)
	DeleteOne(ctx context.Context, id string) (int64, error)
}
