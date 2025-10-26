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

type GetHealth struct {
	verb               string
	path               string
	body               io.Reader
	headers            map[string]string
	ExpectedStatusCode int
	ExpectedErrorMsg   string
	ExpectedBodyMsg    string
}

func TestHealthRoute(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sr),
	)
	tracer := tp.Tracer("test-tracer")

	cases := []GetHealth{
		{
			verb:               "GET",
			path:               "/health",
			body:               nil,
			ExpectedStatusCode: 200,
			ExpectedBodyMsg:    "{\"message\": \"healthy\", \"app\": true, \"database\": true}",
		},
	}

	// var m repository.UserCollectionInterface //nolint:staticcheck
	mh := &mongodb.HealthMockCollectionConfig{
		Error: nil,
	}
	for _, testCase := range cases {

		r, err := NewRouter(tracer, nil, nil, nil, mh)
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

			// In case of expected body as response
			if testCase.ExpectedBodyMsg != "" {
				assert.Equal(t, testCase.ExpectedBodyMsg, rr.Body.String())
			}
		}
	}
}
