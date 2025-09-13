package router

import (
	"budget-tracker-api-v2/internal/http/controller"
	"budget-tracker-api-v2/internal/http/middleware"
	"budget-tracker-api-v2/internal/repository"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewRouter will set new User Routes
func NewRouter(
	userCollectionInterface repository.UserCollectionInterface,
	cardsCollectionInterface repository.CardCollectionInterface,
) (*mux.Router, error) {
	r := mux.NewRouter()

	r.Use(middleware.InjectHeaders)
	r.Use(otelhttp.NewMiddleware("budget-tracker-api-v2"))

	// API routes
	userController := controller.UsersController{
		Repo: userCollectionInterface,
	}

	cardsController := controller.CardsController{
		Repo: cardsCollectionInterface,
	}

	userController.RegisterRoutes(r)
	cardsController.RegisterRoutes(r)

	return r, nil
}
