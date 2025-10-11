package router

import (
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type GetUserTestForAuth struct {
	verb               string
	path               string
	body               io.Reader
	headers            map[string]string
	ExpectedStatusCode int
	ExpectedErrorMsg   string
	ExpectedBodyMsg    string
}

func TestAuthRoute(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sr),
	)
	tracer := tp.Tracer("test-tracer")

	mockedJWTUser := &model.User{
		Login:    "mocked",
		Password: "mockedPassword",
	}

	mockedExistentJWTUser := &model.User{
		Login:    "vsantos",
		Password: "myrandompassword",
	}

	noUserBody, _ := json.Marshal(mockedJWTUser)
	userBody, _ := json.Marshal(mockedExistentJWTUser)

	cases := []GetUserTestForAuth{
		{
			verb:               "POST",
			path:               "/api/v1/jwt/issue",
			body:               nil,
			ExpectedStatusCode: 400,
			ExpectedBodyMsg:    "{\"message\": \"empty body\"}",
		},
		{
			verb:               "POST",
			path:               "/api/v1/jwt/issue",
			body:               bytes.NewBuffer(noUserBody),
			ExpectedStatusCode: 401,
			ExpectedBodyMsg:    "{\"message\": \"invalid credentials\", \"details\": \"unable to find user\"}",
		},
		{
			verb:               "POST",
			path:               "/api/v1/jwt/issue",
			body:               bytes.NewBuffer(userBody),
			ExpectedStatusCode: 201,
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
			req, err := http.NewRequest(testCase.verb, testCase.path, testCase.body)

			// Injecting headers to request
			for hKey, hValue := range testCase.headers {
				req.Header.Add(hKey, hValue)
			}

			assert.NoError(t, err)
			assert.NotNil(t, req)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)
			assert.Equal(t, testCase.ExpectedStatusCode, rr.Code)
			if testCase.ExpectedBodyMsg != "" {
				assert.Equal(t, testCase.ExpectedBodyMsg, rr.Body.String())
			}
		}
	}
}
