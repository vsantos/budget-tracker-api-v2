package main

import (
	"budget-tracker-api-v2/internal/http/router"
	"budget-tracker-api-v2/internal/observability"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	ctx := context.Background()

	timedOut, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	shutdown := observability.InitTracer(timedOut)
	defer func() {
		if cerr := shutdown(timedOut); cerr != nil {
			log.Printf("error shutting down tracer object: %v", cerr)
		}
	}()

	tracer := otel.Tracer("budget-tracker-api-v2")

	c, err := mongodb.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	var m repository.UserCollectionInterface //nolint:staticcheck
	m = &mongodb.UserCollectionConfig{
		Tracer:          tracer,
		MongoCollection: c.Database("budget-tracker-v2").Collection("users"),
	}

	var ms repository.CardCollectionInterface //nolint:staticcheck
	ms = &mongodb.CardCollectionConfig{
		Tracer:          tracer,
		MongoCollection: c.Database("budget-tracker-v2").Collection("cards"),
	}

	var mh repository.HealthCollectionInterface //nolint:staticcheck
	mh = &mongodb.HealthCollectionConfig{
		Tracer:          tracer,
		MongoCollection: c.Database("budget-tracker-v2").Collection("health"),
	}

	router, err := router.NewRouter(tracer, m, ms, mh)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server running on :8080")
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
