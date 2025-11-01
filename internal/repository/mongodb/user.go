package mongodb

import (
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/utils/crypt"
	"fmt"
	"time"

	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
)

// MongoUserRepository defines a Repository for User model
type MongoUserRepository struct {
	MongoCollection repository.UserCollectionInterface
}

// NewUserRepository will return an UserRepoInterface for mongodb
func NewUserRepository(ctx context.Context, tracer trace.Tracer, c repository.UserCollectionInterface) (repository.UserRepoInterface, error) {

	_, span := tracer.Start(ctx, "UsersRepository.NewRepository")
	defer span.End()

	r := MongoUserRepository{
		MongoCollection: c,
	}

	return &r, nil
}

// Insert will insert an user
func (r *MongoUserRepository) Insert(ctx context.Context, emp *model.User) (*model.User, error) {

	if emp.ID.IsZero() {
		emp.ID = primitive.NewObjectID()
	}

	t := time.Now()
	emp.CreatedAt = primitive.NewDateTimeFromTime(t)

	// adding salted password for user
	if emp.Password == "" {
		return &model.User{}, errors.New("empty password input")
	}

	sPassword, err := crypt.GenerateSaltedPassword(emp.Password)
	if err != nil {
		return &model.User{}, err
	}

	emp.Password = sPassword

	rUser, err := r.MongoCollection.
		InsertOne(ctx, emp)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			return nil, errors.New("user or email already registered")
		}

		return nil, err
	}

	return rUser, nil
}

// FindByID will fetch an user based on its ID
func (r *MongoUserRepository) FindByID(ctx context.Context, empID string) (*model.User, error) {
	emp, err := r.MongoCollection.FindOne(ctx, empID)

	if err != nil {
		if strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, fmt.Errorf("user id '%s' not found", empID)
		}

		return nil, err
	}

	return emp, nil
}

// Delete will delete an user based on its ID
func (r *MongoUserRepository) Delete(ctx context.Context, empID string) (int64, error) {
	result, err := r.MongoCollection.
		DeleteOne(context.Background(), empID)

	if err != nil {
		return 0, err
	}

	return result, nil
}
