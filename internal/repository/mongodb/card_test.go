package mongodb

import (
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"context"
	"errors"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type insertCardTest struct {
	collection repository.CardCollectionInterface
	card       *model.Card
	err        string
}

func TestInsertCard(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sr),
	)
	tracer := tp.Tracer("test-tracer")

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
			err:  "card already registered with the 'last 4 digits'",
		},
	}

	for _, test := range inserCardTests {
		u, err := NewCardRepository(context.TODO(), tracer, test.collection)
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
