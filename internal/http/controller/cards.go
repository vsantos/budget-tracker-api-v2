package controller

import (
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

// CardsController injects CardRepository to controllers
type CardsController struct {
	Repo repository.CardCollectionInterface
}

// RegisterRoutes register router for handling Card operations
func (uc *CardsController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/cards", uc.GetCards).Methods("GET")
	r.HandleFunc("/api/v1/cards", uc.CreateCard).Methods("POST")
	r.HandleFunc("/api/v1/cards/{id}", uc.GetCard).Methods("GET")
}

// GetCards handler list of all card within the platform without filters. Deprecated.
func (uc *CardsController) GetCards(w http.ResponseWriter, r *http.Request) {
	var cards []model.Card
	err := json.NewEncoder(w).Encode(cards)
	if err != nil {
		log.Error("Could not encode response: ", err)
	}
}

// CreateCard create a new card within the platform
func (uc *CardsController) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card *model.Card

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	u, err := mongodb.NewCardRepository(context.Background(), uc.Repo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	card, err = u.Insert(r.Context(), card)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			w.WriteHeader(http.StatusConflict)
			_, err := w.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
			if err != nil {
				log.Error("Could not write response: ", err)
			}

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(`{"message": "created card '` + card.Alias + `'", "id": "` + card.ID.Hex() + `", "owner_id": "` + card.OwnerID.Hex() + `"}`))
	if err != nil {
		log.Error("Could not write response: ", err)
	}

}

// GetCard will find a single card based on ID
func (uc *CardsController) GetCard(w http.ResponseWriter, r *http.Request) {
	var card *model.Card
	tracer := otel.Tracer("budget-tracker-api-v2")
	ctx, span := tracer.Start(r.Context(), "controller")
	defer span.End()

	params := mux.Vars(r)

	u, err := mongodb.NewCardRepository(ctx, uc.Repo)
	if err != nil {

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		if strings.Contains(err.Error(), "already registered") {
			w.WriteHeader(http.StatusConflict)
			_, err := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				log.Error("Could not write response: ", err)
			}

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not get card", "details": "` + err.Error() + `"}`))
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			log.Error("Could not write response: ", err)
		}

		span.AddEvent("card info returned")
		return
	}

	card, err = u.FindByID(r.Context(), params["id"])
	if err != nil {
		if strings.Contains(err.Error(), "could not find card") {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": "could not find card", "id": "` + params["id"] + `"}`))
			if err != nil {
				log.Error("Could not write response: ", err)
			}

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "` + err.Error() + `", "owner_id": "` + card.OwnerID.Hex() + `}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	err = json.NewEncoder(w).Encode(card)
	if err != nil {
		log.Error("Could not encode response: ", err)
	}
}
