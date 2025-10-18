package repository

import (
	"budget-tracker-api-v2/internal/model"
	"context"
)

// UserRepoInterface defines User CRUD operations
type UserRepoInterface interface {
	Insert(ctx context.Context, emp *model.User) (*model.User, error)
	FindByID(ctx context.Context, empID string) (*model.User, error)
	Delete(ctx context.Context, id string) (int64, error)
	// FindAllUser(ctx context.Context) ([]model.User, error)
	// UpdateUserByID(ctx context.Context, empID string, updatedEmp *model.User) (int64, error)
	// DeleteUserByID(ctx context.Context, empID string) (int64, error)
}
