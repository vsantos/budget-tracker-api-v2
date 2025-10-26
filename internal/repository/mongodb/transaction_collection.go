package mongodb

import (
	"budget-tracker-api-v2/internal/model"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TransactionCollectionConfig will implement mongodb collection functions
type TransactionCollectionConfig struct {
	Tracer          trace.Tracer
	MongoCollection *mongo.Collection
}

// CreateIndexes will create mongodb indexes
func (c *TransactionCollectionConfig) CreateIndexes(ctx context.Context, indexes []string) error {
	tracer := otel.Tracer("budget-tracker-api-v2")
	indCtx, span := tracer.Start(ctx, "TransactionCollection.CreateIndexes")
	defer span.End()

	var indexModels []mongo.IndexModel
	for _, i := range indexes {
		indexModels = append(indexModels, mongo.IndexModel{
			Keys:    bson.D{{Key: i, Value: 1}},
			Options: options.Index().SetUnique(true),
		})

		kv := attribute.StringSlice("indexes", indexes)
		span.SetAttributes(kv)
	}

	_, err := c.MongoCollection.Indexes().CreateMany(indCtx, indexModels)
	if err != nil {
		return err
	}

	return nil
}

// InsertOne will insert a document into mongodb
func (c *TransactionCollectionConfig) InsertOne(ctx context.Context, t *model.Transaction) (transaction *model.Transaction, err error) {
	ctx, span := c.Tracer.Start(ctx, "TransactionCollection.InsertOne")
	defer span.End()

	r, err := c.MongoCollection.InsertOne(ctx, t)
	fmt.Println(r.InsertedID)
	fmt.Println(err)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return t, nil
}

// FindOne will find a Transaction from collection
func (c *TransactionCollectionConfig) FindOne(ctx context.Context, id string) (*model.Transaction, error) {
	tracer := otel.Tracer("budget-tracker-api-v2")
	fCtx, span := tracer.Start(ctx, "TransactionsCollection.FindOne")
	span.SetAttributes(attribute.String("id", id))
	defer span.End()

	var transaction model.Transaction

	i, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("not able to get ID format")
	}

	filter := bson.M{"_id": i}
	err = c.MongoCollection.FindOne(fCtx, filter).Decode(&transaction)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// DeleteOne will find a Transaction from collection
func (c *TransactionCollectionConfig) DeleteOne(ctx context.Context, id string) (int64, error) {
	ctx, span := c.Tracer.Start(ctx, "TransactionCollection.DeleteOne")
	defer span.End()

	i, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, errors.New("not able to get ID format")
	}

	filter := bson.M{"_id": i}
	result, err := c.MongoCollection.DeleteOne(ctx, filter)

	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
