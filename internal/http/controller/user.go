package controller

import (
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// GetUsers fetches all users from the platform
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	json.NewEncoder(w).Encode(users)
}

// CreateUser fetches all users from the platform
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	var user *model.User

	_ = json.NewDecoder(r.Body).Decode(&user)

	c, err := mongodb.NewClient("mongodb+srv://budget-tracker.gj4ww.mongodb.net")
	if err != nil {
		log.Fatal(err)
		return
	}

	var m repository.UserCollectionInterface
	m = &mongodb.UserCollectionConfig{
		MongoCollection: c.Database("budget-tracker-v2").Collection("users"),
	}

	u, err := mongodb.NewUserRepository(context.Background(), m)
	if err != nil {
		log.Fatal(err)
		return
	}

	user, err = u.Insert(r.Context(), user)
	if err != nil {
		if strings.Contains(err.Error(), "user or email already registered") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"message": "could not create user", "details": "` + err.Error() + `"}`))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "could not create user", "details": "` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "created user '` + user.Login + `'", "id": "` + user.ID.String() + `"}`))
}

// GetUser will find a single user based on ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	json.NewEncoder(w).Encode(user)

	w.Header().Add("content-type", "application/json")
	params := mux.Vars(r)

	c, err := mongodb.NewClient("mongodb+srv://budget-tracker.gj4ww.mongodb.net")
	if err != nil {
		log.Fatal(err)
		return
	}

	var m repository.UserCollectionInterface
	m = &mongodb.UserCollectionConfig{
		MongoCollection: c.Database("budget-tracker-v2").Collection("users"),
	}

	u, err := mongodb.NewUserRepository(context.Background(), m)
	if err != nil {
		log.Fatal(err)
		return
	}

	user, err = u.FindByID(r.Context(), params["id"])
	if err != nil {
		if strings.Contains(err.Error(), "could not find user") {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "could not find user", "id": "` + params["id"] + `"}`))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// masking salted password
	user.Password = ""

	json.NewEncoder(w).Encode(user)
}
