package mongodb

import (
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type insertSpendTest struct {
	collection repository.SpendCollectionInterface
	spend      *model.Spend
	err        string
}

func TestInsertSpend(t *testing.T) {

	var inserSpendTests = []insertSpendTest{
		{
			collection: &SpendMockCollectionConfig{},
			spend:      &model.Spend{},
			err:        "",
		},
		{
			collection: &SpendMockCollectionConfig{},
			spend: &model.Spend{
				ID: primitive.NewObjectID(),
			},
			err: "",
		},
		{
			collection: &SpendMockCollectionConfig{
				Error: errors.New("duplicate key error collection"),
			},
			spend: &model.Spend{},
			err:   "spend already registered with the same ID",
		},
	}

	for _, test := range inserSpendTests {
		u, err := NewSpendRepository(context.TODO(), test.collection)
		assert.NoError(t, err)

		_, err = u.Insert(context.TODO(), test.spend)
		if test.err == "" {
			assert.NoError(t, err)
		}

		if test.err != "" {
			assert.EqualError(t, err, test.err)
		}
	}
}
