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

	_, span := tracer.Start(ctx, "TransactionRepository.NewRepository")
	defer span.End()

	r := MongoTransactionRepository{
		Tracer:          tracer,
		MongoCollection: c,
	}

	// err := r.MongoCollection.CreateIndexes(ctx, []string{"_id"})
	// if err != nil {
	// 	span.RecordError(err)
	// 	span.SetStatus(codes.Error, err.Error())
	// 	return nil, err
	// }
	return &r, nil
}

// Insert will insert an card
func (r *MongoTransactionRepository) Insert(ctx context.Context, emp *model.Transaction) (*model.Transaction, error) {
	ctx, span := r.Tracer.Start(ctx, "TransactionRepository.Insert")
	defer span.End()

	if emp.ID.IsZero() {
		emp.ID = primitive.NewObjectID()
	}

	t := time.Now()
	emp.CreatedAt = primitive.NewDateTimeFromTime(t)

	if emp.TransactionDate == 0 {
		emp.TransactionDate = emp.CreatedAt
	}

	_, err := r.MongoCollection.
		InsertOne(ctx, emp)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			return nil, errors.New("card already registered with the 'last 4 digits'")
		}

		return nil, err
	}

	return emp, nil
}

// FindByID will fetch an card based on its ID
func (r *MongoTransactionRepository) FindByID(ctx context.Context, empID string) (*model.Transaction, error) {
	ctx, span := r.Tracer.Start(ctx, "TransactionRepository.FindByID")
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

// // FindAllTransaction will fetch all card
// func (r *MongoTransactionRepository) FindAllTransaction(ctx context.Context) ([]model.Transaction, error) {
// 	var emps []model.Transaction

// 	results, err := r.MongoCollection.
// 		Find(ctx, bson.D{})

// 	if err != nil {
// 		return nil, err
// 	}

// 	err = results.All(ctx, &emps)
// 	if err != nil {
// 		return nil, errors.New("unable to decode")
// 	}

// 	return emps, nil
// }

// // UpdateTransactionByID will update an card based on its ID
// func (r *MongoTransactionRepository) UpdateTransactionByID(ctx context.Context, empID string, updatedEmp *model.Transaction) (int64, error) {
// 	result, err := r.MongoCollection.
// 		UpdateOne(ctx,
// 			bson.D{
// 				{
// 					Key:   "card_id",
// 					Value: empID,
// 				}},
// 			bson.D{
// 				{
// 					Key:   "$set",
// 					Value: updatedEmp,
// 				}},
// 		)

// 	if err != nil {
// 		return 0, err
// 	}

// 	return result.ModifiedCount, nil
// }

// Delete will delete an card based on its ID
func (r *MongoTransactionRepository) Delete(ctx context.Context, empID string) (int64, error) {
	ctx, span := r.Tracer.Start(ctx, "TransactionRepository.Delete")
	defer span.End()

	result, err := r.MongoCollection.
		DeleteOne(ctx, empID)

	if err != nil {
		return 0, err
	}

	return result, nil
}
