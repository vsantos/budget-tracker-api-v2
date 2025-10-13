package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// HealthCollectionConfig will implement mongodb collection functions for healthchecks
type HealthCollectionConfig struct {
	Tracer          trace.Tracer
	MongoCollection *mongo.Collection
}

func (c *HealthCollectionConfig) Ping(ctx context.Context) (healthy bool, err error) {
	tracer := otel.Tracer("budget-tracker-api-v2")
	spanC, span := tracer.Start(ctx, "healthcheck")
	defer span.End()

	err = c.MongoCollection.Database().Client().Ping(spanC, nil)
	if err != nil {
		return false, err
	}

	return true, nil
}
