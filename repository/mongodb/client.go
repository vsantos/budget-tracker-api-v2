package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient will return a valid mongoDB connection
func NewClient(url string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uri := options.Client().ApplyURI(url)
	cred := options.Client().SetAuth(options.Credential{
		Username: "vsantos",
		Password: "sZgDBYXG2DwxC4qG",
	})

	client, err := mongo.Connect(ctx, uri, cred)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client, nil
}
