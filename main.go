package main

import (
	"budget-tracker-api-v2/model"
	"budget-tracker-api-v2/repository"
	"budget-tracker-api-v2/repository/mongodb"
	"context"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	ctx := context.Background()

	c, err := mongodb.NewClient("mongodb+srv://budget-tracker.gj4ww.mongodb.net")
	if err != nil {
		log.Error(err)
	}

	var m repository.UserCollectionInterface
	m = &mongodb.UserCollectionConfig{
		MongoCollection: c.Database("budget-tracker").Collection("users"),
	}

	u, err := mongodb.NewUserRepository(ctx, m)
	if err != nil {
		log.Error(err)
	}

	r, err := u.InsertUser(ctx, &model.User{
		Login: "vsantos",
		Name:  "Victor Santos",
		Email: "vsantos.py@gmail.com",
	})
	if err != nil {
		log.Error(err)
	}

	if r != nil {
		log.WithFields(log.Fields{
			"id":    r.ID.Hex(),
			"login": r.Login,
		}).Info("User created")

		ru, err := u.FindUserByID(ctx, r.ID.Hex())
		if err != nil {
			log.Error(err)
		}

		if ru != nil {
			log.Info(ru)
		}
	}
}
