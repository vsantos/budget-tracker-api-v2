package controller

import (
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	/*adicionar essa linha */

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// UsersController injects UserRepository to controllers
type UsersController struct {
	Repo repository.UserCollectionInterface
}

// RegisterRoutes register router for handling User operations
func (uc *UsersController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/users", uc.GetUsers).Methods("GET")
	r.HandleFunc("/api/v1/users", uc.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/users/{id}", uc.GetUser).Methods("GET")
}

// GetUsers handler list of all user within the platform without filters. Deprecated.
func (uc *UsersController) GetUsers(w http.ResponseWriter, r *http.Request) {
	// var users []model.User

	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(users)
	_, err := w.Write([]byte(`[]`))
	if err != nil {
		log.Error("Could not write response: ", err)
	}
}

// CreateUser create a new user within the platform
func (uc *UsersController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user *model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create user", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}
		return
	}

	u, err := mongodb.NewUserRepository(context.Background(), uc.Repo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create user", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	user, err = u.Insert(r.Context(), user)
	if err != nil {
		if strings.Contains(err.Error(), "user or email already registered") {
			w.WriteHeader(http.StatusConflict)
			_, err := w.Write([]byte(`{"message": "could not create user", "details": "` + err.Error() + `"}`))
			if err != nil {
				log.Error("Could not write response: ", err)
			}

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create user", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(`{"message": "created user '` + user.Login + `'", "id": "` + user.ID.Hex() + `"}`))
	if err != nil {
		log.Error("Could not write response: ", err)
	}

}

// GetUser will find a single user based on ID
func (uc *UsersController) GetUser(w http.ResponseWriter, r *http.Request) {
	var user *model.User

	params := mux.Vars(r)

	u, err := mongodb.NewUserRepository(context.Background(), uc.Repo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not get user", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	user, err = u.FindByID(r.Context(), params["id"])
	if err != nil {
		if strings.Contains(err.Error(), "could not find user") {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": "could not find user", "id": "` + params["id"] + `"}`))
			if err != nil {
				log.Error("Could not write response: ", err)
			}

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	// masking salted password
	user.Password = "<sensitive>"

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Error("Could not encode response: ", err)
	}
}
