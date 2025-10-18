package controller

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func TestGenerateJWTAccessToken(t *testing.T) {
	var jwtKey = "myhellokey"

	token, err := GenerateJWTAccessToken(t.Context(), jwtKey, "foo", "bar")
	assert.NoError(t, err)
	t.Log(token)
	to, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("non expected signature: %v", t.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "HS256", to.Header["alg"])
	assert.Equal(t, "JWT", to.Header["typ"])
}
