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

	ru, err := u.FindByID(ctx, r.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	if ru == nil {
		log.Panic(ru)
	}

	var mc repository.CardCollectionInterface
	mc = &mongodb.CardCollectionConfig{
		MongoCollection: c.Database("budget-tracker").Collection("cards"),
	}

	sc, err := mongodb.NewCardRepository(ctx, mc)
	if err != nil {
		log.Fatal(err)
	}

	cardOutput, err := sc.Insert(ctx, &model.Card{
		OwnerID:    ru.ID,
		Alias:      "foo",
		Network:    "VISA",
		LastDigits: 5443,
	})
	if err != nil {
		log.Fatal(err)
	}

	rc, err := sc.FindByID(ctx, cardOutput.ID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	if ru == nil {
		log.Panic(ru)
	}

	_, err = s.Insert(ctx, &model.Spend{
		OwnerID:     ru.ID,
		Type:        "variant",
		Description: "Aluguel do mês",
		Cost:        30.2,
		Categories:  []string{"casa"},
		PaymentMethod: model.PaymentMethod{
			Credit:      *rc,
			Debit:       false,
			PaymentSlip: false,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
