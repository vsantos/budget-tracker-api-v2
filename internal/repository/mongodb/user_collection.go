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
	"go.opentelemetry.io/otel/trace"
)

// UserCollectionConfig will implement mongodb collection functions
type UserCollectionConfig struct {
	Tracer          trace.Tracer
	MongoCollection *mongo.Collection
}

// CreateIndexes will create mongodb indexes
func (c *UserCollectionConfig) CreateIndexes(ctx context.Context, indexes []string) error {
	tracer := otel.Tracer("budget-tracker-api-v2")
	fCtx, span := tracer.Start(ctx, "insert one")
	defer span.End()

	var indexModels []mongo.IndexModel
	for _, i := range indexes {
		indexModels = append(indexModels, mongo.IndexModel{
			Keys:    bson.D{{Key: i, Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	}

	_, err := c.MongoCollection.Indexes().CreateMany(fCtx, indexModels)
	if err != nil {
		return err
	}

	return nil
}

// InsertOne will insert a document into mongodb
func (c *UserCollectionConfig) InsertOne(ctx context.Context, document interface{}) (id string, err error) {
	tracer := otel.Tracer("budget-tracker-api-v2")
	fCtx, span := tracer.Start(ctx, "insert one")
	defer span.End()

	r, err := c.MongoCollection.InsertOne(fCtx, document)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", r.InsertedID), nil
}

// FindOne will find a User from collection
func (c *UserCollectionConfig) FindOne(ctx context.Context, id string) (*model.User, error) {
	var emp model.User

	tracer := otel.Tracer("budget-tracker-api-v2")
	fCtx, span := tracer.Start(ctx, "find one")
	defer span.End()

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

// FindOne will find a User from collection
func (c *UserCollectionConfig) FindOneBy(ctx context.Context, login string) (*model.User, error) {
	var emp model.User

	tracer := otel.Tracer("budget-tracker-api-v2")
	fCtx, span := tracer.Start(ctx, "find one by")
	defer span.End()

	filter := bson.M{"login": login}
	err := c.MongoCollection.FindOne(fCtx, filter).Decode(&emp)
	if err != nil {
		return nil, err
	}

	// span.SetAttributes()

	return &emp, nil
}

// DeleteOne will find a User from collection
func (c *UserCollectionConfig) DeleteOne(ctx context.Context, id string) (int64, error) {
	tracer := otel.Tracer("budget-tracker-api-v2")
	fCtx, span := tracer.Start(ctx, "delete one")
	defer span.End()

	i, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, errors.New("not able to get ID format")
	}

	filter := bson.M{"_id": i}
	result, err := c.MongoCollection.DeleteOne(fCtx, filter)

	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
