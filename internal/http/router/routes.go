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

	r.Handle("/swagger/swagger.yaml", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))
	// r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swaggerui"))))
	r.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `
	<!DOCTYPE html>
	<html>
	<head>
	  <title>Swagger UI</title>
	  <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui.css" />
	</head>
	<body>
	  <div id="swagger-ui"></div>
	  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui-bundle.js"></script>
	  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.18.3/swagger-ui-standalone-preset.js"></script>
	  <script>
	    window.onload = function() {
	      SwaggerUIBundle({
	        url: "/swagger/swagger.yaml",
	        dom_id: '#swagger-ui',
	        presets: [
	          SwaggerUIBundle.presets.apis,
	          SwaggerUIStandalonePreset
	        ],
	        layout: "BaseLayout"
	      });
	    }
	  </script>
	</body>
	</html>
			`
		w.Write([]byte(html))
	})
	// API routes
	userController := controller.UsersController{
		Tracer: tracer,
		Repo:   userCollectionInterface,
	}

	cardsController := controller.CardsController{
		Tracer: tracer,
		Repo:   cardsCollectionInterface,
	}

	userController.RegisterRoutes(r)
	cardsController.RegisterRoutes(r)

	return r, nil
}
