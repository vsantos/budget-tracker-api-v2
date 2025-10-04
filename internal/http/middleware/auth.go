package middleware

import (
	"fmt"
	"mime"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

var mySigninKey = []byte("myhellokey")

// RequireContentTypeJSON enforces JSON content-type from requests
func RequireContentTypeJSON(h http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*")

		contentType := request.Header.Get("Content-Type")

		if contentType == "" {
			response.WriteHeader(http.StatusBadRequest)
			_, err := response.Write([]byte(`{"message": "empty Content-Type header"}`))
			if err != nil {
				log.Error("Could not write response: ", err)
			}
			return
		}
		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				response.WriteHeader(http.StatusBadRequest)
				_, err := response.Write([]byte(`{"message": "malformed Content-Type header"}`))
				if err != nil {
					log.Error("Could not write response: ", err)
				}
				return
			}

			if mt != "application/json" {
				response.WriteHeader(http.StatusUnsupportedMediaType)
				_, err := response.Write([]byte(`{"message": "content-Type header must be application/json"}`))
				if err != nil {
					log.Error("Could not write response: ", err)
				}
				return
			}
		}

		h.ServeHTTP(response, request)
	})
}

// RequireTokenAuthentication enforces authentication token from requests
func RequireTokenAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*")

		if request.Header["Authorization"] == nil {
			response.WriteHeader(http.StatusBadRequest)
			_, err := response.Write([]byte(`{"message": "missing 'Authorization' header"}`))
			if err != nil {
				log.Error("Could not write response: ", err)
			}
			return
		}

		if request.Header["Authorization"] != nil {
			jwtString := strings.Split(request.Header["Authorization"][0], "Bearer ")
			if len(jwtString) <= 1 {
				response.WriteHeader(http.StatusUnauthorized)
				_, err := response.Write([]byte(`{"message": "could not parse token", "details": "possible mistyped bearer token"}`))
				if err != nil {
					log.Error("Could not write response: ", err)
				}
				return
			}

			token, err := jwt.Parse(jwtString[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("could not decode token")
				}
				return mySigninKey, nil
			})

			if err != nil {
				response.WriteHeader(http.StatusUnauthorized)
				_, err := response.Write([]byte(`{"message": "could not authenticate", "details": "` + err.Error() + `"}`))
				if err != nil {
					log.Error("Could not write response: ", err)
				}
				return
			}

			if !token.Valid {
				response.WriteHeader(http.StatusInternalServerError)
				_, err := response.Write([]byte(`{"message": "token not valid"}`))
				if err != nil {
					log.Error("Could not write response: ", err)
				}
			}
		}

		h.ServeHTTP(response, request)
	})
}
