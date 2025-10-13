package mongodb

import "budget-tracker-api-v2/internal/repository"

// MongoHealthRepository defines a Repository for User model
type MongoHealthRepository struct {
	MongoCollection repository.HealthCollectionInterface
}
