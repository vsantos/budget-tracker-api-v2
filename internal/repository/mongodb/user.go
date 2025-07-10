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
)

// MongoUserRepository defines a Repository for User model
type MongoUserRepository struct {
	MongoCollection repository.UserCollectionInterface
}

// NewUserRepository will return an UserRepoInterface for mongodb
func NewUserRepository(ctx context.Context, c repository.UserCollectionInterface) (repository.UserRepoInterface, error) {
	r := MongoUserRepository{
		MongoCollection: c,
	}

	err := r.MongoCollection.CreateIndexes(context.TODO(), []string{"login", "email"})
	if err != nil {
		return nil, err
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

	_, err = r.MongoCollection.
		InsertOne(context.Background(), emp)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error collection") {
			return nil, errors.New("user or email already registered")
		}

		return nil, err
	}

	return emp, nil
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

// // FindAllUser will fetch all user
// func (r *MongoUserRepository) FindAllUser(ctx context.Context) ([]model.User, error) {
// 	var emps []model.User

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

// // UpdateUserByID will update an user based on its ID
// func (r *MongoUserRepository) UpdateUserByID(ctx context.Context, empID string, updatedEmp *model.User) (int64, error) {
// 	result, err := r.MongoCollection.
// 		UpdateOne(context.Background(),
// 			bson.D{
// 				{
// 					Key:   "user_id",
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

// Delete will delete an user based on its ID
func (r *MongoUserRepository) Delete(ctx context.Context, empID string) (int64, error) {
	result, err := r.MongoCollection.
		DeleteOne(context.Background(), empID)

	if err != nil {
		return 0, err
	}

	return result, nil
}
