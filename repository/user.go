package repository

import (
	"budget-tracker-api-v2/model"
	"context"
)

// UserRepoInterface defines User CRUD operations
type UserRepoInterface interface {
	InsertUser(ctx context.Context, emp *model.User) (*model.User, error)
	FindUserByID(ctx context.Context, empID string) (*model.User, error)
	// FindAllUser(ctx context.Context) ([]model.User, error)
	// UpdateUserByID(ctx context.Context, empID string, updatedEmp *model.User) (int64, error)
	// DeleteUserByID(ctx context.Context, empID string) (int64, error)
}
