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
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type TransactionTestCase struct {
	verb               string
	path               string
	body               io.Reader
	headers            map[string]string
	ExpectedStatusCode int
	ExpectedErrorMsg   string
	ExpectedBodyMsg    string
}

func TestGetTransactionRoute(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sr),
	)
	tracer := tp.Tracer("test-tracer")

	invalidJwtToken, err := controller.GenerateJWTAccessToken(t.Context(), "invalidKey", "ID_MALUCO", "vsantos")
	assert.NoError(t, err)

	validJwtToken, err := controller.GenerateJWTAccessToken(t.Context(), "myhellokey", "ID_MALUCO", "vsantos")
	assert.NoError(t, err)

	cases := []TransactionTestCase{
		{
			verb:               "GET",
			path:               "/api/v1/transactions",
			body:               nil,
			ExpectedStatusCode: 405,
			ExpectedErrorMsg:   "",
		},
		{
			verb:               "GET",
			path:               "/api/v1/transactionssss",
			body:               nil,
			ExpectedStatusCode: 404,
			ExpectedErrorMsg:   "",
		},
		{
			verb: "GET",
			path: "/api/v1/transactions/68fd6b00f4c9e77e59aaf97e",
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
			path: "/api/v1/transactions/68fd6b00f4c9e77e59aaf97e",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", invalidJwtToken),
			},
			ExpectedStatusCode: 401,
			ExpectedBodyMsg:    "{\"message\": \"could not authenticate\", \"details\": \"signature is invalid\"}",
		},
		{
			verb: "GET",
			path: "/api/v1/transactions/68fd6b00f4c9e77e59aaf97e",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", validJwtToken),
			},
			ExpectedStatusCode: 200,
			ExpectedBodyMsg:    "{\"id\":\"68fd6b00f4c9e77e59aaf97e\",\"balance_id\":\"687baad049572fb8c4e305f3\",\"owner_id\":\"66f1cca3c37c733c4ada103d\",\"type\":\"income\",\"description\":\"My favorite chinese restaurant\",\"amount\":15.3,\"currency\":\"BRL\",\"payment_method\":{\"credit\":{\"id\":\"000000000000000000000000\",\"owner_id\":\"000000000000000000000000\",\"alias\":\"\",\"type\":\"\",\"network\":\"\",\"bank\":\"\",\"last_digits\":0},\"pix\":false,\"payment_slip\":false},\"transaction_date\":\"2023-10-26T14:30:00Z\",\"categories\":[\"food\"],\"created_at\":\"2023-10-26T14:30:00Z\"}\n",
		},
		{
			verb: "GET",
			path: "/api/v1/transactions/686f255205535b1dd3b68f39",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", validJwtToken),
			},
			ExpectedStatusCode: 404,
			ExpectedBodyMsg:    "{\"message\": \"could not find transaction\", \"id\": \"686f255205535b1dd3b68f39\"}",
		},
		{
			verb: "GET",
			path: "/api/v1/transactions/686f255205535b1dd3b68f3ss",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", validJwtToken),
			},
			ExpectedStatusCode: 500,
			ExpectedBodyMsg:    "{\"message\": \"the provided hex string is not a valid ObjectID\"}",
		},
	}

	// var m repository.UserCollectionInterface //nolint:staticcheck
	um := &mongodb.UserMockCollectionConfig{
		Error: nil,
	}
	tm := &mongodb.TransactionMockCollectionConfig{
		Error: nil,
	}
	cm := &mongodb.CardMockCollectionConfig{
		Error: nil,
	}

	for _, testCase := range cases {

		r, err := NewRouter(tracer, um, cm, tm, nil)
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
			t.Log(testCase.path)
			t.Log(rr.Body.String())

			// In case of expected body as response
			if testCase.ExpectedBodyMsg != "" {
				assert.Equal(t, testCase.ExpectedBodyMsg, rr.Body.String())
			}
		}
	}
}

func TestCreateTransactionRoute(t *testing.T) {
	objID, err := primitive.ObjectIDFromHex("68fd6b00f4c9e77e59aaf97e")
	assert.NoError(t, err)
	balanceID, err := primitive.ObjectIDFromHex("687baad049572fb8c4e305f3")
	assert.NoError(t, err)
	ownerID, err := primitive.ObjectIDFromHex("66f1cca3c37c733c4ada103d")
	assert.NoError(t, err)

	mockedTransactionDate := "2023-10-26 14:30:00"
	layout := "2006-01-02 15:04:05"
	tDate, err := time.Parse(layout, mockedTransactionDate)
	assert.NoError(t, err)

	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sr),
	)
	tracer := tp.Tracer("test-tracer")

	invalidJwtToken, err := controller.GenerateJWTAccessToken(t.Context(), "invalidKey", "ID_MALUCO", "vsantos")
	assert.NoError(t, err)

	validJwtToken, err := controller.GenerateJWTAccessToken(t.Context(), "myhellokey", "ID_MALUCO", "vsantos")
	assert.NoError(t, err)

	mockedSuccessBody := &model.Transaction{
		ID:              objID,
		BalanceID:       balanceID,
		OwnerID:         ownerID,
		Type:            "income",
		Description:     "My favorite chinese restaurant",
		Amount:          15.3,
		Currency:        "BRL",
		PaymentMethod:   model.PaymentMethod{},
		TransactionDate: primitive.NewDateTimeFromTime(tDate),
		Categories:      []string{"food"},
		CreatedAt:       primitive.NewDateTimeFromTime(tDate),
	}

	// Encode struct to JSON
	mockedJSONBody, err := json.Marshal(mockedSuccessBody)
	assert.NoError(t, err)

	cases := []TransactionTestCase{
		{
			verb:               "POST",
			path:               "/api/v1/transactions",
			body:               nil,
			ExpectedStatusCode: 400,
		},
		{
			verb:               "POST",
			path:               "/api/v1/transactions/68fd6b00f4c9e77e59aaf97e",
			body:               nil,
			ExpectedStatusCode: 404,
		},
		{
			verb: "POST",
			path: "/api/v1/transactions",
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
			path: "/api/v1/transactions",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", invalidJwtToken),
			},
			ExpectedStatusCode: 401,
			ExpectedBodyMsg:    "{\"message\": \"could not authenticate\", \"details\": \"signature is invalid\"}",
		},
		{
			verb: "POST",
			path: "/api/v1/transactions",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", validJwtToken),
			},
			ExpectedStatusCode: 400,
			ExpectedBodyMsg:    "{\"message\": \"could not create transaction\", \"details\": \"missing body\"}",
		},
		{
			verb: "POST",
			path: "/api/v1/transactions",
			headers: map[string]string{
				"Content-type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", validJwtToken),
			},
			ExpectedStatusCode: 201,
			body:               bytes.NewBuffer(mockedJSONBody),
			ExpectedBodyMsg:    "{\"message\":\"transaction created\",\"id\":\"68fd6b00f4c9e77e59aaf97e\",\"owner_id\":\"66f1cca3c37c733c4ada103d\",\"status_code\":201,\"transaction\":{\"id\":\"68fd6b00f4c9e77e59aaf97e\",\"balance_id\":\"687baad049572fb8c4e305f3\",\"owner_id\":\"66f1cca3c37c733c4ada103d\",\"type\":\"income\",\"description\":\"My favorite chinese restaurant\",\"amount\":15.3,\"currency\":\"BRL\",\"payment_method\":{\"credit\":{\"id\":\"000000000000000000000000\",\"owner_id\":\"000000000000000000000000\",\"alias\":\"\",\"type\":\"\",\"network\":\"\",\"bank\":\"\",\"last_digits\":0},\"pix\":false,\"payment_slip\":false},\"transaction_date\":\"2023-10-26T14:30:00Z\",\"categories\":[\"food\"],\"created_at\":\"2023-10-26T14:30:00Z\"}}\n",
		},
	}

	// var m repository.UserCollectionInterface //nolint:staticcheck
	um := &mongodb.UserMockCollectionConfig{
		Error: nil,
	}
	tm := &mongodb.TransactionMockCollectionConfig{
		Error: nil,
	}
	cm := &mongodb.CardMockCollectionConfig{
		Error: nil,
	}

	for _, testCase := range cases {

		r, err := NewRouter(tracer, um, cm, tm, nil)
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
