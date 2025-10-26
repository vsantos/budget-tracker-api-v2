package repository

import (
	"budget-tracker-api-v2/internal/model"
	"context"
)

// TransactionRepoInterface defines Card CRUD operations
type TransactionRepoInterface interface {
	Insert(ctx context.Context, emp *model.Transaction) (*model.Transaction, error)
	FindByID(ctx context.Context, empID string) (*model.Transaction, error)
	Delete(ctx context.Context, id string) (int64, error)
	// FindAllCard(ctx context.Context) ([]model.Card, error)
	// UpdateCardByID(ctx context.Context, empID string, updatedEmp *model.Card) (int64, error)
	// DeleteCardByID(ctx context.Context, empID string) (int64, error)
}
