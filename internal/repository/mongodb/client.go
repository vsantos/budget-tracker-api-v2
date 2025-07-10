package mongodb

import (
	"context"
	"errors"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient will return a valid mongoDB connection
func NewClient(url string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var mongoUser, mongoPass string
	mongoUser = os.Getenv("MONGODB_ATLAS_USER")
	mongoPass = os.Getenv("MONGODB_ATLAS_PASS")

	if mongoUser == "" || mongoPass == "" {
		return nil, errors.New("empty USER or PASS env vars for mongodb atlas")
	}

	uri := options.Client().ApplyURI(url)
	cred := options.Client().SetAuth(options.Credential{
		Username: mongoUser,
		Password: mongoPass,
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
