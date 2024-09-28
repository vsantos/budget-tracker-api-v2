package mongodb

import (
	"budget-tracker-api-v2/model"
	"budget-tracker-api-v2/repository"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type insertUserTest struct {
	collection repository.UserCollectionInterface
	user       *model.User
	err        string
}

func TestInsertUser(t *testing.T) {

	var inserUserTests = []insertUserTest{
		{
			collection: &UserMockCollectionConfig{},
			user:       &model.User{},
			err:        "",
		},
		{
			collection: &UserMockCollectionConfig{},
			user: &model.User{
				ID: primitive.NewObjectID(),
			},
			err: "",
		},
		{
			collection: &UserMockCollectionConfig{
				Error: errors.New("duplicate key error collection"),
			},
			user: &model.User{},
			err:  "user or email already registered",
		},
	}

	for _, test := range inserUserTests {
		u, err := NewUserRepository(context.TODO(), test.collection)
		assert.NoError(t, err)

		_, err = u.Insert(context.TODO(), test.user)
		if test.err == "" {
			assert.NoError(t, err)
		}

		if test.err != "" {
			assert.EqualError(t, err, test.err)
		}
	}
}
