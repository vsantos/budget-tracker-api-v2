package router

import (
	"budget-tracker-api-v2/internal/http/controller"
	"budget-tracker-api-v2/internal/http/middleware"
	"budget-tracker-api-v2/internal/repository"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"github.com/gorilla/mux"
)

// NewRouter will set new User Routes
func NewRouter(
	userCollectionInterface repository.UserCollectionInterface,
	cardsCollectionInterface repository.CardCollectionInterface,
) (*mux.Router, error) {
	r := mux.NewRouter()

	r.Use(middleware.InjectHeaders)
	r.Use(
		otelmux.Middleware(
			"budget-tracker-api-v2",
			otelmux.WithSpanNameFormatter(func(routeName string, r *http.Request) string {
				return fmt.Sprintf("%s %s", r.Method, routeName)
			}),
		),
	)
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
