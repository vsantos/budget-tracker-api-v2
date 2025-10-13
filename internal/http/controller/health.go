package controller

import (
	"budget-tracker-api-v2/internal/repository"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type HealthController struct {
	Tracer     trace.Tracer
	HealthRepo repository.HealthCollectionInterface
}

// RegisterRoutes register router for handling healthcheck operations
func (uc *HealthController) RegisterRoutes(r *mux.Router) {
	p := r.PathPrefix("/health").Subrouter()

	p.HandleFunc("", uc.HealthCheck).Methods("GET")
}

// Ping handler list of all card within the platform without filters. Deprecated.
func (uc *HealthController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	sCtx, span := uc.Tracer.Start(r.Context(), "CardsController.Health")
	defer span.End()

	status, err := uc.HealthRepo.Ping(sCtx)
	if !status || err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "unhealthy", "app": true, "database": false}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "healthy", "app": true, "database": true}`))
	if err != nil {
		log.Error("Could not write response: ", err)
	}
}
