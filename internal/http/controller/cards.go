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
	json.NewEncoder(w).Encode(cards)
}

// CreateCard create a new card within the platform
func (uc *CardsController) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card *model.Card

	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
		return
	}

	u, err := mongodb.NewCardRepository(context.Background(), uc.Repo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
		return
	}

	card, err = u.Insert(r.Context(), card)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "could not create card", "details": "` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "created card '` + card.Alias + `'", "id": "` + card.ID.Hex() + `", "owner_id": "` + card.OwnerID.Hex() + `"}`))
}

// GetCard will find a single card based on ID
func (uc *CardsController) GetCard(w http.ResponseWriter, r *http.Request) {
	var card *model.Card

	params := mux.Vars(r)

	u, err := mongodb.NewCardRepository(context.Background(), uc.Repo)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"message": "` + err.Error() + `"}`))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "could not get card", "details": "` + err.Error() + `"}`))
		return
	}

	card, err = u.FindByID(r.Context(), params["id"])
	if err != nil {
		if strings.Contains(err.Error(), "could not find card") {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "could not find card", "id": "` + params["id"] + `"}`))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "` + err.Error() + `", "owner_id": "` + card.OwnerID.Hex() + `}`))
		return
	}

	json.NewEncoder(w).Encode(card)
}
