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

type insertCardTest struct {
	collection repository.CardCollectionInterface
	card       *model.Card
	err        string
}

func TestInsertCard(t *testing.T) {

	var inserCardTests = []insertCardTest{
		{
			collection: &CardMockCollectionConfig{},
			card:       &model.Card{},
			err:        "",
		},
		{
			collection: &CardMockCollectionConfig{},
			card: &model.Card{
				ID: primitive.NewObjectID(),
			},
			err: "",
		},
		{
			collection: &CardMockCollectionConfig{
				Error: errors.New("duplicate key error collection"),
			},
			card: &model.Card{},
			err:  "card already registered with the same ID and/or owner ID",
		},
	}

	for _, test := range inserCardTests {
		u, err := NewCardRepository(context.TODO(), test.collection)
		assert.NoError(t, err)

		_, err = u.Insert(context.TODO(), test.card)
		if test.err == "" {
			assert.NoError(t, err)
		}

		if test.err != "" {
			assert.EqualError(t, err, test.err)
		}
	}
}
