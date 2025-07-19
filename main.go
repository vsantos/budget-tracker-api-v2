package main

import (
	"budget-tracker-api-v2/internal/http/router"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	// ctx := context.Background()

	// c, err := mongodb.NewClient("mongodb+srv://budget-tracker.gj4ww.mongodb.net")
	c, err := mongodb.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	var m repository.UserCollectionInterface
	m = &mongodb.UserCollectionConfig{
		MongoCollection: c.Database("budget-tracker-v2").Collection("users"),
	}

	var ms repository.CardCollectionInterface
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
