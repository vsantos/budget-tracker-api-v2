package mongodb

import (
	"budget-tracker-api-v2/model"
	"budget-tracker-api-v2/repository"
	"fmt"
	"time"

	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoSpendRepository defines a Repository for Spend model
type MongoSpendRepository struct {
	MongoCollection repository.SpendCollectionInterface
}

// NewSpendRepository will return an SpendRepoInterface for mongodb
func NewSpendRepository(ctx context.Context, c repository.SpendCollectionInterface) (repository.SpendRepoInterface, error) {
	r := MongoSpendRepository{
		MongoCollection: c,
	}

	// err := r.MongoCollection.CreateIndexes(context.TODO(), []string{"login", "email"})
	// if err != nil {
	// 	return nil, err
	// }
	return &r, nil
}

// Insert will insert an spend
func (r *MongoSpendRepository) Insert(ctx context.Context, emp *model.Spend) (*model.Spend, error) {

	if emp.ID.IsZero() {
		emp.ID = primitive.NewObjectID()
	}

	t := time.Now()
	emp.CreatedAt = primitive.NewDateTimeFromTime(t)

	_, err := r.MongoCollection.
		InsertOne(context.Background(), emp)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			return nil, errors.New("spend already registered with the same ID")
		}

		return nil, err
	}

	return emp, nil
}

// FindByID will fetch an spend based on its ID
func (r *MongoSpendRepository) FindByID(ctx context.Context, empID string) (*model.Spend, error) {
	emp, err := r.MongoCollection.FindOne(ctx, empID)

	if err != nil {
		if strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, fmt.Errorf("spend id '%s' not found", empID)
		}

		return nil, err
	}

	return emp, nil
}

// // FindAllSpend will fetch all spend
// func (r *MongoSpendRepository) FindAllSpend(ctx context.Context) ([]model.Spend, error) {
// 	var emps []model.Spend

// 	results, err := r.MongoCollection.
// 		Find(context.Background(), bson.D{})

// 	if err != nil {
// 		return nil, err
// 	}

// 	err = results.All(context.Background(), &emps)
// 	if err != nil {
// 		return nil, errors.New("unable to decode")
// 	}

// 	return emps, nil
// }

// // UpdateSpendByID will update an spend based on its ID
// func (r *MongoSpendRepository) UpdateSpendByID(ctx context.Context, empID string, updatedEmp *model.Spend) (int64, error) {
// 	result, err := r.MongoCollection.
// 		UpdateOne(context.Background(),
// 			bson.D{
// 				{
// 					Key:   "spend_id",
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

// Delete will delete an spend based on its ID
func (r *MongoSpendRepository) Delete(ctx context.Context, empID string) (int64, error) {
	result, err := r.MongoCollection.
		DeleteOne(context.Background(), empID)

	if err != nil {
		return 0, err
	}

	return result, nil
}
