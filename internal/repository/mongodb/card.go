package mongodb

import (
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"fmt"
	"time"

	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoCardRepository defines a Repository for Card model
type MongoCardRepository struct {
	MongoCollection repository.CardCollectionInterface
}

// NewCardRepository will return an CardRepoInterface for mongodb
func NewCardRepository(ctx context.Context, c repository.CardCollectionInterface) (repository.CardRepoInterface, error) {
	r := MongoCardRepository{
		MongoCollection: c,
	}

	err := r.MongoCollection.CreateIndexes(context.TODO(), []string{"owner_id", "last_digits"})
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Insert will insert an card
func (r *MongoCardRepository) Insert(ctx context.Context, emp *model.Card) (*model.Card, error) {

	if emp.ID.IsZero() {
		emp.ID = primitive.NewObjectID()
	}

	t := time.Now()
	emp.CreatedAt = primitive.NewDateTimeFromTime(t)

	_, err := r.MongoCollection.
		InsertOne(context.Background(), emp)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			return nil, errors.New("card already registered with the same ID and/or owner ID")
		}

		return nil, err
	}

	return emp, nil
}

// FindByID will fetch an card based on its ID
func (r *MongoCardRepository) FindByID(ctx context.Context, empID string) (*model.Card, error) {
	emp, err := r.MongoCollection.FindOne(ctx, empID)

	if err != nil {
		if strings.Contains(err.Error(), "mongo: no documents in result") {
			return nil, fmt.Errorf("card id '%s' not found", empID)
		}

		return nil, err
	}

	return emp, nil
}

// // FindAllCard will fetch all card
// func (r *MongoCardRepository) FindAllCard(ctx context.Context) ([]model.Card, error) {
// 	var emps []model.Card

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

// // UpdateCardByID will update an card based on its ID
// func (r *MongoCardRepository) UpdateCardByID(ctx context.Context, empID string, updatedEmp *model.Card) (int64, error) {
// 	result, err := r.MongoCollection.
// 		UpdateOne(context.Background(),
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
func (r *MongoCardRepository) Delete(ctx context.Context, empID string) (int64, error) {
	result, err := r.MongoCollection.
		DeleteOne(context.Background(), empID)

	if err != nil {
		return 0, err
	}

	return result, nil
}
