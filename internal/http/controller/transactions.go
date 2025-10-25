package controller

import (
	"budget-tracker-api-v2/internal/http/middleware"
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TransactionsController injects CardRepository to controllers
type TransactionsController struct {
	Tracer          trace.Tracer
	TransactionRepo repository.TransactionCollectionInterface
	CardsRepo       repository.CardCollectionInterface
}

type TransactionErrorMessage struct {
	Message    string `json:"message"`
	Details    string `json:"details"`
	StatusCode int32  `json:"status_code,omitempty"`
}

type TransactionCreatedMessage struct {
	Message     string            `json:"message"`
	ID          string            `json:"id"`
	OwnerID     string            `json:"owner_id"`
	StatusCode  int32             `json:"status_code"`
	Transaction model.Transaction `json:"transaction"`
}

// RegisterRoutes register router for handling Card operations
func (uc *TransactionsController) RegisterRoutes(r *mux.Router) {
	p := r.PathPrefix("/api/v1/transactions").Subrouter()
	p.Use(middleware.RequireTokenAuthentication)

	p.HandleFunc("", uc.CreateTransaction).Methods("POST")
}

func (tc *TransactionsController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction *model.Transaction

	ctx, span := tc.Tracer.Start(r.Context(), "TransactionsController.CreateTransaction")
	defer span.End()

	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create transaction", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	t, err := mongodb.NewTransactionRepository(ctx, tc.Tracer, tc.TransactionRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create transaction", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	c, err := mongodb.NewCardRepository(ctx, tc.Tracer, tc.CardsRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create transaction", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	cardFilter := bson.M{"alias": transaction.PaymentMethod.Credit.Alias}
	card, err := c.FindByFilter(ctx, cardFilter)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())

		if strings.Contains(err.Error(), "not found") {
			span.AddEvent("card not found")

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

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create transaction", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}
	transaction.PaymentMethod.Credit = *card

	rt, err := t.Insert(ctx, transaction)
	if err != nil {

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		if strings.Contains(err.Error(), "already registered") {
			w.WriteHeader(http.StatusConflict)
			msg := TransactionErrorMessage{
				Message:    "could not create transaction",
				Details:    err.Error(),
				StatusCode: http.StatusConflict,
			}
			if err := json.NewEncoder(w).Encode(msg); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Error("could not encode response for transaction creation", err)
				return
			}

			log.Error(msg)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create transaction", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("could not write response for transaction creation", err)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	sMsg := TransactionCreatedMessage{
		Message:     "transaction created",
		ID:          rt.ID.Hex(),
		OwnerID:     rt.OwnerID.Hex(),
		StatusCode:  http.StatusCreated,
		Transaction: *rt,
	}
	if err := json.NewEncoder(w).Encode(sMsg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("could not encode response for transaction creation", err)
		return
	}

	log.Info(sMsg)
}
