package controller

import (
	"budget-tracker-api-v2/internal/http/middleware"
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Due to a limitation on Swagger 2.0, we are declaring static error messages
type CardsErrorMessage struct {
	Message    string `json:"message"`
	Details    string `json:"details"`
	StatusCode int32  `json:"status_code,omitempty"`
}

type CardDeletedMessage struct {
	Message    string `json:"message"`
	ID         string `json:"id"`
	StatusCode int32  `json:"status_code,omitempty"`
}
type CardsCreatedMessage struct {
	Message    string     `json:"message"`
	ID         string     `json:"id"`
	OwnerID    string     `json:"owner_id"`
	StatusCode int32      `json:"status_code"`
	Card       model.Card `json:"card"`
}

// CardsController injects CardRepository to controllers
type CardsController struct {
	Tracer trace.Tracer
	Repo   repository.CardCollectionInterface
}

// RegisterRoutes register router for handling Card operations
func (uc *CardsController) RegisterRoutes(r *mux.Router) {
	p := r.PathPrefix("/api/v1/cards").Subrouter()
	p.Use(middleware.RequireTokenAuthentication)

	p.HandleFunc("", uc.GetCards).Methods("GET")
	p.HandleFunc("", uc.CreateCard).Methods("POST")
	p.HandleFunc("/{id}", uc.GetCard).Methods("GET")
	p.HandleFunc("/{id}", uc.DeleteCard).Methods("DELETE")
}

// GetCards handler list of all card within the platform without filters. Deprecated.
func (uc *CardsController) GetCards(w http.ResponseWriter, r *http.Request) {
	_, span := uc.Tracer.Start(r.Context(), "CardsController.GetCards")
	defer span.End()

	var cards []model.Card
	err := json.NewEncoder(w).Encode(cards)
	if err != nil {
		log.Error("Could not encode response: ", err)
	}
}

func (uc *CardsController) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card *model.Card

	ctx, span := uc.Tracer.Start(r.Context(), "CardsController.CreateCard")
	defer span.End()

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	u, err := mongodb.NewCardRepository(ctx, uc.Tracer, uc.Repo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	card, err = u.Insert(ctx, card)
	if err != nil {

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		if strings.Contains(err.Error(), "already registered") {
			w.WriteHeader(http.StatusConflict)
			msg := CardsErrorMessage{
				Message:    "could not create card",
				Details:    err.Error(),
				StatusCode: http.StatusConflict,
			}
			if err := json.NewEncoder(w).Encode(msg); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
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
	sMsg := CardsCreatedMessage{
		Message:    "card created",
		ID:         card.ID.Hex(),
		OwnerID:    card.OwnerID.Hex(),
		StatusCode: http.StatusCreated,
		Card:       *card,
	}
	if err := json.NewEncoder(w).Encode(sMsg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (uc *CardsController) GetCard(w http.ResponseWriter, r *http.Request) {
	var card *model.Card

	ctx, span := uc.Tracer.Start(r.Context(), "CardsController.GetCard")
	defer span.End()

	params := mux.Vars(r)

	u, err := mongodb.NewCardRepository(ctx, uc.Tracer, uc.Repo)
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

	card, err = u.FindByID(ctx, params["id"])
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		if strings.Contains(err.Error(), "not found") {
			span.AddEvent("user not found")

			notFoundMsg := CardsErrorMessage{
				Message:    "could not find card",
				Details:    err.Error(),
				StatusCode: http.StatusNotFound,
			}

			w.WriteHeader(http.StatusNotFound)
			err := json.NewEncoder(w).Encode(notFoundMsg)
			if err != nil {
				log.Error("Could not write response: ", err)
			}

			return
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		errMsg := CardsErrorMessage{
			Message: "error when fetching card's details",
			Details: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(errMsg)
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

func (uc *CardsController) DeleteCard(w http.ResponseWriter, r *http.Request) {
	ctx, span := uc.Tracer.Start(r.Context(), "CardsController.DeleteCard")
	defer span.End()

	params := mux.Vars(r)
	resp, err := uc.Repo.DeleteOne(ctx, params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		erroMsg := CardsErrorMessage{
			Message:    "could not find card",
			Details:    err.Error(),
			StatusCode: http.StatusInternalServerError,
		}

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		err = json.NewEncoder(w).Encode(erroMsg)
		if err != nil {
			log.Error("Could not encode response: ", err)
		}
		return
	}

	if resp == 0 {
		w.WriteHeader(http.StatusNotFound)
		errMsg := CardsErrorMessage{
			Message:    "card not deleted",
			Details:    fmt.Sprintf("no cards were deleted from given card ID '%s'", params["id"]),
			StatusCode: http.StatusNotFound,
		}
		err = json.NewEncoder(w).Encode(errMsg)
		if err != nil {
			log.Error("Could not write response: ", err)
		}
		span.AddEvent("no cards deleted")
		return
	}

	delMsg := CardDeletedMessage{
		Message:    "card deleted",
		ID:         params["id"],
		StatusCode: http.StatusOK,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(delMsg)
	if err != nil {
		log.Error("Could not write response: ", err)
	}

	span.AddEvent("card deleted")
}
