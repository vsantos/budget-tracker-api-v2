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
func NewClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var mongoHost, mongoUser, mongoPass string

	mongoHost = os.Getenv("MONGODB_HOST")
	mongoUser = os.Getenv("MONGODB_USER")
	mongoPass = os.Getenv("MONGODB_PASS")

	if mongoHost == "" || mongoUser == "" || mongoPass == "" {
		return nil, errors.New("empty HOST, USER or PASS env vars for mongodb atlas")
	}

	uri := options.Client().ApplyURI(mongoHost)
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
