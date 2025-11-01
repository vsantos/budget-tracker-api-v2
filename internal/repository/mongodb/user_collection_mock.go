package mongodb

import (
	"budget-tracker-api-v2/internal/model"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	usertDate time.Time
)

func init() {
	mockedTransactionDate := "2023-10-26 14:30:00"
	layout := "2006-01-02 15:04:05"
	usertDate, _ = time.Parse(layout, mockedTransactionDate)
}

// UserMockCollectionConfig will implement mongodb collection functions
type UserMockCollectionConfig struct {
	Error error
}

// CreateIndexes will create mongodb indexes
func (c *UserMockCollectionConfig) CreateIndexes(ctx context.Context, indexes []string) error {
	return nil
}

// InsertOne will insert a document into mongodb
func (c *UserMockCollectionConfig) InsertOne(ctx context.Context, emp *model.User) (user *model.User, err error) {
	objID, err := primitive.ObjectIDFromHex("686f255205535b1dd3b68f38")
	if err != nil {
		return nil, err
	}

	if c.Error != nil {
		return nil, c.Error
	}

	fmt.Println("mocked inserted")
	fmt.Println(objID)
	return &model.User{
		ID:        objID,
		Login:     emp.Login,
		Firstname: emp.Firstname,
		Lastname:  emp.Lastname,
		Email:     emp.Email,
		Password:  "$2a$10$HOrmuqyfwr575K4P9tjQXe0QKqbddMA/KFZ.YZhWVKPLMUF3LS4gi",
		CreatedAt: primitive.NewDateTimeFromTime(usertDate),
	}, nil
}

// FindOneBy will find a document based on field
func (c *UserMockCollectionConfig) FindOneBy(ctx context.Context, login string) (*model.User, error) {
	objID, err := primitive.ObjectIDFromHex("686f255205535b1dd3b68f38")
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        objID,
		Login:     login,
		Firstname: "Mock User",
		Lastname:  "Torres",
		Email:     "mock.user@gmail.com",
		Password:  "$2a$10$HOrmuqyfwr575K4P9tjQXe0QKqbddMA/KFZ.YZhWVKPLMUF3LS4gi",
	}, nil
}

// FindOne will find a document based on ID
func (c *UserMockCollectionConfig) FindOne(ctx context.Context, id string) (*model.User, error) {
	objID, err := primitive.ObjectIDFromHex("686f255205535b1dd3b68f38")
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        objID,
		Login:     "mockuser",
		Firstname: "Mock User",
		Lastname:  "Torres",
		Email:     "mock.user@gmail.com",
		// Salted password
		Password: "$2a$10$HOrmuqyfwr575K4P9tjQXe0QKqbddMA/KFZ.YZhWVKPLMUF3LS4gi",
	}, nil
}

// DeleteOne will insert a document into mongodb
func (c *UserMockCollectionConfig) DeleteOne(ctx context.Context, id string) (int64, error) {
	return 1, nil
}
