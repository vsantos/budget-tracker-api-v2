package repository

import (
	"budget-tracker-api-v2/model"
	"context"
)

// CardRepoInterface defines Card CRUD operations
type CardRepoInterface interface {
	Insert(ctx context.Context, emp *model.Card) (*model.Card, error)
	FindByID(ctx context.Context, empID string) (*model.Card, error)
	Delete(ctx context.Context, id string) (int64, error)
	// FindAllCard(ctx context.Context) ([]model.Card, error)
	// UpdateCardByID(ctx context.Context, empID string, updatedEmp *model.Card) (int64, error)
	// DeleteCardByID(ctx context.Context, empID string) (int64, error)
}
