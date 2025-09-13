package main

import (
	"budget-tracker-api-v2/internal/http/router"
	"budget-tracker-api-v2/internal/obsevability"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	ctx := context.Background()

	timedOut, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	shutdown := obsevability.InitTracer(timedOut)
	defer shutdown(timedOut)

	c, err := mongodb.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	var m repository.UserCollectionInterface //nolint:staticcheck
	m = &mongodb.UserCollectionConfig{
		MongoCollection: c.Database("budget-tracker-v2").Collection("users"),
	}

	var ms repository.CardCollectionInterface //nolint:staticcheck
	ms = &mongodb.CardCollectionConfig{
		MongoCollection: c.Database("budget-tracker-v2").Collection("cards"),
	}

	router, err := router.NewRouter(m, ms)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
