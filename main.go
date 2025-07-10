package main

import (
	"budget-tracker-api-v2/internal/http/router"
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		MongoCollection: c.Database("budget-tracker-v2").Collection("users"),
	}

	u, err := mongodb.NewUserRepository(ctx, m)
	if err != nil {
		log.Fatal(err)
	}

	r, err := u.Insert(ctx, &model.User{
		Login:     "vsantos",
		Firstname: "Victor",
		Lastname:  "Santos",
		Email:     "vsantos.py@gmail.com",
		Password:  "MySuperSecretPassword",
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

	userID, _ := primitive.ObjectIDFromHex("686f255205535b1dd3b68f38")
	ru, err := u.FindByID(ctx, userID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	if ru == nil || ru.ID.IsZero() {
		log.Panic(ru)
	}

	var mc repository.CardCollectionInterface
	mc = &mongodb.CardCollectionConfig{
		MongoCollection: c.Database("budget-tracker-v2").Collection("cards"),
	}

	sc, err := mongodb.NewCardRepository(ctx, mc)
	if err != nil {
		log.Fatal(err)
	}

	cardOutput, err := sc.Insert(ctx, &model.Card{
		OwnerID:    ru.ID,
		Alias:      "platinum multiplo",
		Network:    "VISA",
		Bank:       "Itaú",
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

	var ms repository.SpendCollectionInterface
	ms = &mongodb.SpendCollectionConfig{
		MongoCollection: c.Database("budget-tracker-v2").Collection("spends"),
	}

	s, err := mongodb.NewSpendRepository(ctx, ms)
	if err != nil {
		log.Fatal(err)
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

	router := router.NewRouter()
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
