package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInjectHeaders(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := InjectHeaders(dummyHandler)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
}
