package router

import (
	"budget-tracker-api-v2/internal/repository"
	"budget-tracker-api-v2/internal/repository/mongodb"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type GetUserTest struct {
	verb               string
	path               string
	body               io.Reader
	ExpectedStatusCode int
	ExpectedErrorMsg   string
}

func TestGetUserRoute(t *testing.T) {
	t.Setenv("MONGODB_ATLAS_USER", "user")
	t.Setenv("MONGODB_ATLAS_PASS", "pass")
	t.Setenv("MONGODB_ATLAS_HOST", "mongodb+srv://budget-tracker.gj4ww.mongodb.net")

	cases := []GetUserTest{
		{
			verb:               "GET",
			path:               "/api/v1/users",
			body:               nil,
			ExpectedStatusCode: 200,
			ExpectedErrorMsg:   "",
		},
		{
			verb:               "GET",
			path:               "/api/v1/users",
			body:               nil,
			ExpectedStatusCode: 200,
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
			verb:               "GET",
			path:               "/api/v1/users/686f255205535b1dd3b68f38",
			body:               nil,
			ExpectedStatusCode: 200,
			ExpectedErrorMsg:   "",
		},
	}

	var m repository.UserCollectionInterface
	m = &mongodb.UserMockCollectionConfig{
		Error: nil,
	}

	for _, testCase := range cases {

		r, err := NewRouter(m, nil)
		if testCase.ExpectedErrorMsg != "" {
			assert.Error(t, err, testCase.ExpectedErrorMsg)
		} else {
			assert.NoError(t, err)
		}

		if err == nil {
			req, err := http.NewRequest(testCase.verb, testCase.path, nil)
			assert.NoError(t, err)
			assert.NotNil(t, req)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, testCase.ExpectedStatusCode)
		}
	}

}
