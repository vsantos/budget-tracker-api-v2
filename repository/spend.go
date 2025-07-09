package repository

import (
	"budget-tracker-api-v2/model"
	"context"
)

// SpendRepoInterface defines Spend CRUD operations
type SpendRepoInterface interface {
	Insert(ctx context.Context, emp *model.Spend) (*model.Spend, error)
	FindByID(ctx context.Context, empID string) (*model.Spend, error)
	Delete(ctx context.Context, id string) (int64, error)
	// FindAllSpend(ctx context.Context) ([]model.Spend, error)
	// UpdateSpendByID(ctx context.Context, empID string, updatedEmp *model.Spend) (int64, error)
	// DeleteSpendByID(ctx context.Context, empID string) (int64, error)
}
