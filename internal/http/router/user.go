package router

import (
	"budget-tracker-api-v2/internal/http/controller"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/v1/users", controller.GetUsers).Methods("GET")
	r.HandleFunc("/api/v1/users", controller.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/users/{id}", controller.GetUser).Methods("GET")

	return r
}
