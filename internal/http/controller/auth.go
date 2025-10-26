package controller

import (
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/utils/crypt"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type AuthController struct {
	Tracer   trace.Tracer
	UserRepo repository.UserCollectionInterface
}

// GenerateJWTAccessToken will generate a JWT access token
func GenerateJWTAccessToken(ctx context.Context, jwtKey string, sub string, login string) (string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["sub"] = sub
	claims["name"] = login
	claims["exp"] = time.Now().Add(5 * time.Minute).Unix()
	claims["iat"] = time.Now().Unix()

	at, err := accessToken.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return at, nil
}

// GenerateJWTRefreshToken will generate a new refresh token
func GenerateJWTRefreshToken(ctx context.Context, jwtKey string, sub string) (string, error) {
	bjwtKey := []byte(jwtKey)
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = sub
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString(bjwtKey)
	if err != nil {
		return "", err
	}

	return rt, nil
}

// RegisterRoutes register router for handling Card operations
func (ac *AuthController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/jwt/issue", ac.CreateToken).Methods("POST")
}

// GetCards handler list of all card within the platform without filters. Deprecated.
func (uc *AuthController) CreateToken(w http.ResponseWriter, r *http.Request) {
	ctx, span := uc.Tracer.Start(r.Context(), "AuthController.CreateToken")
	defer span.End()

	var jwtUser model.JWTUser

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "could not create token", "details": "missing body"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}

		return
	}

	err := json.NewDecoder(r.Body).Decode(&jwtUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Could not write response: ", err)
	}

	if jwtUser.Login == "" || jwtUser.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"message": "required login and password to request a new token"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}
		return
	}

	user, err := uc.UserRepo.FindOneBy(ctx, jwtUser.Login)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	var aToken, rToken string

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"message": "could not find given user"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}
		return
	}

	match := crypt.CheckPasswordHash(jwtUser.Password, user.Password)
	if match {
		aToken, err = GenerateJWTAccessToken(ctx, "myhellokey", user.ID.Hex(), user.Login)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(`{"message": "could not create access token", "details": "` + err.Error() + `"}`))
			if err != nil {
				log.Error("Could not write response: ", err)
			}
			return
		}

		rToken, err = GenerateJWTRefreshToken(ctx, "myhellokey", user.ID.Hex())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(`{"message": "could not create refresh token", "details": "` + err.Error() + `"}`))
			if err != nil {
				log.Error("Could not write response: ", err)
			}
			return
		}

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		err = errors.New("invalid credentials")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		_, err := w.Write([]byte(`{"message": "invalid credentials", "details": "unable to find user"}`))
		if err != nil {
			log.Error("Could not write response: ", err)
		}
		return
	}
	// }

	// w.WriteHeader(http.StatusUnauthorized)
	jwtResponse := &model.JWTResponse{
		Type:         "bearer",
		RefreshToken: rToken,
		AccessToken:  aToken,
	}

	jwtResponseJSON, err := json.Marshal(jwtResponse)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("Could not write response: ", err)
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(jwtResponseJSON)
	if err != nil {
		log.Error("Could not write response: ", err)
	}
}
