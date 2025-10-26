package repository

import (
	"budget-tracker-api-v2/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// CardRepoInterface defines Card CRUD operations
type CardRepoInterface interface {
	Insert(ctx context.Context, emp *model.Card) (*model.Card, error)
	FindByID(ctx context.Context, empID string) (*model.Card, error)
	FindByFilter(ctx context.Context, filter bson.M) (*model.Card, error)
	Delete(ctx context.Context, id string) (int64, error)
}
