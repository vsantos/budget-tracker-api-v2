package repository

import (
	"context"
)

// UserCollectionInterface defines a mongodb collection API to be posteriorly mocked
type HealthCollectionInterface interface {
	Ping(ctx context.Context) (bool, error)
}
