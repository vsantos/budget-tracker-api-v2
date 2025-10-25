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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// CardCollectionConfig will implement mongodb collection functions
type CardCollectionConfig struct {
	Tracer          trace.Tracer
	MongoCollection *mongo.Collection
}

// CreateIndexes will create mongodb indexes
func (c *CardCollectionConfig) CreateIndexes(ctx context.Context, indexes []string) error {
	tracer := otel.Tracer("budget-tracker-api-v2")
	indCtx, span := tracer.Start(ctx, "create indexes")
	defer span.End()

	var indexModels []mongo.IndexModel
	for _, i := range indexes {
		indexModels = append(indexModels, mongo.IndexModel{
			Keys:    bson.D{{Key: i, Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	}

	_, err := c.MongoCollection.Indexes().CreateMany(indCtx, indexModels)
	if err != nil {
		return err
	}

	return nil
}

// InsertOne will insert a document into mongodb
func (c *CardCollectionConfig) InsertOne(ctx context.Context, document interface{}) (id string, err error) {
	ctx, span := c.Tracer.Start(ctx, "CardCollection.InsertOne")
	defer span.End()

	r, err := c.MongoCollection.InsertOne(ctx, document)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", err
	}

	return fmt.Sprintf("%v", r.InsertedID), nil
}

// FindOne will find a Card from collection
func (c *CardCollectionConfig) FindOne(ctx context.Context, id string) (*model.Card, error) {
	tracer := otel.Tracer("budget-tracker-api-v2")
	fCtx, span := tracer.Start(ctx, "find one")
	defer span.End()

	var emp model.Card

	i, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("not able to get ID format")
	}

	filter := bson.M{"_id": i}
	err = c.MongoCollection.FindOne(fCtx, filter).Decode(&emp)

	if err != nil {
		return nil, err
	}

	return &emp, nil
}

// FindOne will find a Card from collection
func (c *CardCollectionConfig) FindOneByFilter(ctx context.Context, filter bson.M) (*model.Card, error) {
	tracer := otel.Tracer("budget-tracker-api-v2")
	fCtx, span := tracer.Start(ctx, "find one")
	defer span.End()

	var emp model.Card

	err := c.MongoCollection.FindOne(fCtx, filter).Decode(&emp)
	if err != nil {
		return nil, err
	}

	return &emp, nil
}

// DeleteOne will find a Card from collection
func (c *CardCollectionConfig) DeleteOne(ctx context.Context, id string) (int64, error) {
	ctx, span := c.Tracer.Start(ctx, "CardCollection.DeleteOne")
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
