package router

import (
	"budget-tracker-api-v2/internal/http/controller"
	"budget-tracker-api-v2/internal/model"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"bytes"
	"encoding/json"
	"fmt"
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

	invalidJwtToken, err := controller.GenerateJWTAccessToken(t.Context(), "invalidKey", "ID_MALUCO", "vsantos")
	assert.NoError(t, err)

	validJwtToken, err := controller.GenerateJWTAccessToken(t.Context(), "myhellokey", "ID_MALUCO", "vsantos")
	assert.NoError(t, err)

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
			ExpectedBodyMsg:    "{\"message\": \"could not authenticate\", \"details\": \"Token is expired\"}",
		},
		{
			verb: "GET",
			path: "/api/v1/users/686f255205535b1dd3b68f38",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", invalidJwtToken),
			},
			ExpectedStatusCode: 401,
			ExpectedBodyMsg:    "{\"message\": \"could not authenticate\", \"details\": \"signature is invalid\"}",
		},
		{
			verb: "GET",
			path: "/api/v1/users/686f255205535b1dd3b68f38",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", validJwtToken),
			},
			ExpectedStatusCode: 200,
			ExpectedBodyMsg:    "{\"id\":\"686f255205535b1dd3b68f38\",\"login\":\"mockuser\",\"firstname\":\"Mock User\",\"lastname\":\"Torres\",\"email\":\"mock.user@gmail.com\",\"password\":\"\\u003csensitive\\u003e\"}\n",
		},
	}

	// var m repository.UserCollectionInterface //nolint:staticcheck
	m := &mongodb.UserMockCollectionConfig{
		Error: nil,
	}
	for _, testCase := range cases {

		r, err := NewRouter(tracer, m, nil, nil, nil)
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

func TestCreateUserRoute(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sr),
	)
	tracer := tp.Tracer("test-tracer")

	invalidJwtToken, err := controller.GenerateJWTAccessToken(t.Context(), "invalidKey", "ID_MALUCO", "vsantos")
	assert.NoError(t, err)

	validJwtToken, err := controller.GenerateJWTAccessToken(t.Context(), "myhellokey", "ID_MALUCO", "vsantos")
	assert.NoError(t, err)

	mockedSuccessBody := &model.User{
		Login:     "mocker_user",
		Firstname: "My Mock",
		Lastname:  "Harrinson",
		Email:     "mock@domain.io",
		Password:  "plaintext_pass",
	}
	mockedJSONBody, err := json.Marshal(mockedSuccessBody)
	assert.NoError(t, err)

	cases := []GetUserTest{
		{
			verb:               "POST",
			path:               "/api/v1/userssss",
			body:               nil,
			ExpectedStatusCode: 404,
			ExpectedErrorMsg:   "",
		},
		{
			verb: "POST",
			path: "/api/v1/users",
			body: nil,
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NTk1NDM1NjgsImlhdCI6MTc1OTU0MzI2OCwibmFtZSI6InZzYW50b3MiLCJzdWIiOiI2OGQwNjZlOGE0YjYyYjgwNTE3N2FhOTEifQ.OoL_fl_qz0bmNvFxAh_PoG5UJsDtGCCzeGMa9sPKqVs",
			},
			ExpectedStatusCode: 401,
			ExpectedBodyMsg:    "{\"message\": \"could not authenticate\", \"details\": \"Token is expired\"}",
		},
		{
			verb: "POST",
			path: "/api/v1/users",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", invalidJwtToken),
			},
			ExpectedStatusCode: 401,
			ExpectedBodyMsg:    "{\"message\": \"could not authenticate\", \"details\": \"signature is invalid\"}",
		},
		{
			verb: "POST",
			path: "/api/v1/users",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", validJwtToken),
			},
			body:               nil,
			ExpectedStatusCode: 400,
			ExpectedBodyMsg:    "{\"message\": \"could not create user\", \"details\": \"missing body\"}",
		},
		{
			verb: "POST",
			path: "/api/v1/users",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", validJwtToken),
			},
			body:               bytes.NewBuffer(mockedJSONBody),
			ExpectedStatusCode: 201,
			ExpectedBodyMsg:    "{\"message\": \"created user 'mocker_user'\", \"id\": \"686f255205535b1dd3b68f38\"}",
		},
	}

	// var m repository.UserCollectionInterface //nolint:staticcheck
	m := &mongodb.UserMockCollectionConfig{
		Error: nil,
	}
	for _, testCase := range cases {

		r, err := NewRouter(tracer, m, nil, nil, nil)
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

			// In case of expected body as response
			if testCase.ExpectedBodyMsg != "" {
				assert.Equal(t, testCase.ExpectedBodyMsg, rr.Body.String())
			}
		}
	}
}
