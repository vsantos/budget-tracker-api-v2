package mongodb

import (
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
)

type MongoTransactionRepository struct {
	Tracer          trace.Tracer
	MongoCollection repository.TransactionCollectionInterface
}

// NewTransactionRepository will return an TransactionRepoInterface for mongodb
func NewTransactionRepository(ctx context.Context, tracer trace.Tracer, c repository.TransactionCollectionInterface) (repository.TransactionRepoInterface, error) {

	_, span := tracer.Start(ctx, "TransactionsRepository.NewRepository")
	defer span.End()

	r := MongoTransactionRepository{
		Tracer:          tracer,
		MongoCollection: c,
	}

	return &r, nil
}

// Insert will insert an card
func (r *MongoTransactionRepository) Insert(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	ctx, span := r.Tracer.Start(ctx, "TransactionsRepository.Insert")
	defer span.End()

	if transaction.ID.IsZero() {
		transaction.ID = primitive.NewObjectID()
	}

	t := time.Now()
	transaction.CreatedAt = primitive.NewDateTimeFromTime(t)

	if transaction.TransactionDate == 0 {
		transaction.TransactionDate = transaction.CreatedAt
	}

	returnedTransaction, err := r.MongoCollection.
		InsertOne(ctx, transaction)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			return nil, errors.New("card already registered with the 'last 4 digits'")
		}

		return nil, err
	}

	return returnedTransaction, nil
}

// FindByID will fetch an card based on its ID
func (r *MongoTransactionRepository) FindByID(ctx context.Context, empID string) (*model.Transaction, error) {
	ctx, span := r.Tracer.Start(ctx, "TransactionsRepository.FindByID")
	defer span.End()

	emp, err := r.MongoCollection.FindOne(ctx, empID)

	if err != nil {
		if strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, fmt.Errorf("card id '%s' not found", empID)
		}

		return nil, err
	}

	return emp, nil
}

// Delete will delete an card based on its ID
func (r *MongoTransactionRepository) Delete(ctx context.Context, empID string) (int64, error) {
	ctx, span := r.Tracer.Start(ctx, "TransactionsRepository.Delete")
	defer span.End()

	result, err := r.MongoCollection.
		DeleteOne(ctx, empID)

	if err != nil {
		return 0, err
	}

	return result, nil
}
