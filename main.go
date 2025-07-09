package main

import (
	"budget-tracker-api-v2/model"
	"budget-tracker-api-v2/repository"
	"budget-tracker-api-v2/repository/mongodb"
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	ctx := context.Background()

	c, err := mongodb.NewClient("mongodb+srv://budget-tracker.gj4ww.mongodb.net")
	if err != nil {
		log.Fatal(err)
	}

	var m repository.UserCollectionInterface
	m = &mongodb.UserCollectionConfig{
		MongoCollection: c.Database("budget-tracker").Collection("users"),
	}

	u, err := mongodb.NewUserRepository(ctx, m)
	if err != nil {
		log.Fatal(err)
	}

	r, err := u.Insert(ctx, &model.User{
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

		// _, err = u.Delete(ctx, r.ID.Hex())
		// if err != nil {
		// 	log.Fatal(err)
		// }

	}

	var ms repository.SpendCollectionInterface
	ms = &mongodb.SpendCollectionConfig{
		MongoCollection: c.Database("budget-tracker").Collection("spends"),
	}

	s, err := mongodb.NewSpendRepository(ctx, ms)
	if err != nil {
		log.Fatal(err)
	}

	ru, err := u.FindByID(ctx, "686ec8c798948f5b4911eb68")
	if err != nil {
		log.Fatal(err)
	}

	if ru == nil {
		log.Panic(ru)
	}

	rs, err := s.Insert(ctx, &model.Spend{
		OwnerID:     ru.ID,
		Type:        "variant",
		Description: "Aluguel do mês",
		Cost:        30.2,
		Categories:  []string{"casa"},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(&rs)

	// log.Info(rd)
}
