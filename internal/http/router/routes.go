package router

import (
	"budget-tracker-api-v2/internal/http/controller"
	"budget-tracker-api-v2/internal/http/middleware"
	"budget-tracker-api-v2/internal/repository"
	"context"
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
	transactionsCollectionInterface repository.TransactionCollectionInterface,
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

	var err error

	// API Routes
	userController := controller.UsersController{
		Tracer: tracer,
		Repo:   userCollectionInterface,
	}

	cardsController := controller.CardsController{
		Tracer: tracer,
		Repo:   cardsCollectionInterface,
	}

	transactionsController := controller.TransactionsController{
		Tracer:          tracer,
		TransactionRepo: transactionsCollectionInterface,
		CardsRepo:       cardsCollectionInterface,
	}

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

	fmt.Println("creating indexes")
	if userController.Repo != nil {
		fmt.Println("not null")
		err = userController.Repo.CreateIndexes(context.Background(), []string{"login", "email"})
		if err != nil {
			return nil, err
		}
	}

	cardsController.RegisterRoutes(r)

	if cardsController.Repo != nil {
		err = cardsController.Repo.CreateIndexes(context.Background(), []string{"last_digits"})
		if err != nil {
			return nil, err
		}
	}

	transactionsController.RegisterRoutes(r)
	healthController.RegisterRoutes(r)

	return r, nil
}
