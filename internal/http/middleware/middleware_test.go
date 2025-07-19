package middleware

// import (
// 	"budget-tracker-api-v2/internal/http/router"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestGetUserRoute(t *testing.T) {
// 	type middlewareTestCase struct {
// 		Database      string
// 		ExpectedError error
// 	}

// 	cases := []middlewareTestCase{
// 		{
// 			Database:      "mongodbd",
// 			ExpectedError: nil,
// 		},
// 	}

// 	for _, testCase := range cases {
// 		r, err := router.NewRouter(testCase.Database)
// 		if testCase.ExpectedError != nil {
// 			assert.Error(t, err, testCase.ExpectedError)
// 			assert.NotNil(t, r)
// 		} else {
// 			assert.NoError(t, err)
// 		}
// 	}

// }
