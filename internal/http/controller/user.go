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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// UsersController injects UserRepository to controllers
type UsersController struct {
	Tracer trace.Tracer
	Repo   repository.UserCollectionInterface
}

// RegisterRoutes register router for handling User operations
func (uc *UsersController) RegisterRoutes(r *mux.Router) {
	p := r.PathPrefix("/api/v1/users").Subrouter()
	p.Use(middleware.RequireTokenAuthentication)

	p.HandleFunc("", uc.CreateUser).Methods("POST")
	p.HandleFunc("/{id}", uc.GetUser).Methods("GET")
}

// CreateUser create a new user within the platform
func (uc *UsersController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user *model.User

	ctx, span := uc.Tracer.Start(r.Context(), "UsersController.CreateUser")
	defer span.End()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not create user", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}
		return
	}

	u, err := mongodb.NewUserRepository(ctx, uc.Tracer, uc.Repo)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

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

	ctx, span := uc.Tracer.Start(r.Context(), "UsersController.GetUser")
	defer span.End()

	params := mux.Vars(r)

	u, err := mongodb.NewUserRepository(ctx, uc.Tracer, uc.Repo)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"message": "could not get user", "details": "` + err.Error() + `"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	user, err = u.FindByID(r.Context(), params["id"])
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			span.AddEvent("user not found")

			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte(`{"message": "could not find user", "id": "` + params["id"] + `"}`))
			if err != nil {
				log.Error("Could not write response: ", err)
			}

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

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
