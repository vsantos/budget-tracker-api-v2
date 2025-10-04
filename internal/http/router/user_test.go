package router

import (
	"budget-tracker-api-v2/internal/repository/mongodb"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type GetUserTest struct {
	verb               string
	path               string
	body               io.Reader
	headers            map[string]string
	ExpectedStatusCode int
	ExpectedErrorMsg   string
	ExpectedBodyMsg    string
}

func TestGetUserRoute(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sr),
	)
	tracer := tp.Tracer("test-tracer")

	t.Setenv("MONGODB_ATLAS_USER", "user")
	t.Setenv("MONGODB_ATLAS_PASS", "pass")
	t.Setenv("MONGODB_ATLAS_HOST", "mongodb+srv://budget-tracker.gj4ww.mongodb.net")

	cases := []GetUserTest{
		{
			verb:               "GET",
			path:               "/api/v1/users",
			body:               nil,
			ExpectedStatusCode: 404,
			ExpectedErrorMsg:   "",
		},
		{
			verb:               "GET",
			path:               "/api/v1/userssss",
			body:               nil,
			ExpectedStatusCode: 404,
			ExpectedErrorMsg:   "",
		},
		{
			verb: "GET",
			path: "/api/v1/users/686f255205535b1dd3b68f38",
			body: nil,
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NTk1NDM1NjgsImlhdCI6MTc1OTU0MzI2OCwibmFtZSI6InZzYW50b3MiLCJzdWIiOiI2OGQwNjZlOGE0YjYyYjgwNTE3N2FhOTEifQ.OoL_fl_qz0bmNvFxAh_PoG5UJsDtGCCzeGMa9sPKqVs",
			},
			ExpectedStatusCode: 401,
			ExpectedErrorMsg:   "",
			ExpectedBodyMsg:    "{\"message\": \"could not authenticate\", \"details\": \"Token is expired\"}",
		},
	}

	// var m repository.UserCollectionInterface //nolint:staticcheck
	m := &mongodb.UserMockCollectionConfig{
		Error: nil,
	}
	for _, testCase := range cases {

		r, err := NewRouter(tracer, m, nil)
		if testCase.ExpectedErrorMsg != "" {
			assert.Error(t, err, testCase.ExpectedErrorMsg)
		} else {
			assert.NoError(t, err)
		}

		if err == nil {
			req, err := http.NewRequest(testCase.verb, testCase.path, nil)

			// Injecting headers to request
			for hKey, hValue := range testCase.headers {
				req.Header.Add(hKey, hValue)
			}

			assert.NoError(t, err)
			assert.NotNil(t, req)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)
			assert.Equal(t, testCase.ExpectedStatusCode, rr.Code)
			// t.Log(rr.Body)
			// t.Log(rr.Code)
			// t.Log(testCase)
			if testCase.ExpectedBodyMsg != "" {
				assert.Equal(t, testCase.ExpectedBodyMsg, rr.Body.String())
			}
		}
	}
}
