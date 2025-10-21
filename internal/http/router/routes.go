package router

import (
	"budget-tracker-api-v2/internal/http/controller"
	"budget-tracker-api-v2/internal/http/middleware"
	"budget-tracker-api-v2/internal/repository"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/trace"
)

// NewRouter will set new User Routes
func NewRouter(
	tracer trace.Tracer,
	userCollectionInterface repository.UserCollectionInterface,
	cardsCollectionInterface repository.CardCollectionInterface,
<<<<<<< HEAD
	transactionsCollectionInterface repository.TransactionCollectionInterface,
=======
>>>>>>> main
	healthCollectionInterface repository.HealthCollectionInterface,
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

	// API Routes
	userController := controller.UsersController{
		Tracer: tracer,
		Repo:   userCollectionInterface,
	}

	cardsController := controller.CardsController{
		Tracer: tracer,
		Repo:   cardsCollectionInterface,
	}

<<<<<<< HEAD
	transactionsController := controller.TransactionsController{
		Tracer: tracer,
		Repo:   transactionsCollectionInterface,
	}

=======
>>>>>>> main
	authController := controller.AuthController{
		Tracer:   tracer,
		UserRepo: userController.Repo,
	}

	healthController := controller.HealthController{
		Tracer:     tracer,
		HealthRepo: healthCollectionInterface,
	}

	controller.SwaggerRegisterRouter(r)
	authController.RegisterRoutes(r)
	userController.RegisterRoutes(r)
	cardsController.RegisterRoutes(r)
	transactionsController.RegisterRoutes(r)
	healthController.RegisterRoutes(r)

	return r, nil
}
