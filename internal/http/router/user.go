package router

import (
	"budget-tracker-api-v2/internal/http/controller"
	"budget-tracker-api-v2/internal/http/middleware"
	"budget-tracker-api-v2/internal/repository"

	"github.com/gorilla/mux"
)

// NewRouter will set new User Routes
func NewRouter(userCollectionInterface repository.UserCollectionInterface) (*mux.Router, error) {
	r := mux.NewRouter()
	controller := controller.UsersController{
		Repo: nil,
	}

	controller.Repo = userCollectionInterface

	// API routes
	r.Use(middleware.InjectHeaders)
	r.HandleFunc("/api/v1/users", controller.GetUsers).Methods("GET")
	r.HandleFunc("/api/v1/users", controller.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/users/{id}", controller.GetUser).Methods("GET")

	return r, nil
}
