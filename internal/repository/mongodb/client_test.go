package mongodb

import (
	"io"
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

func TestClientNoEnvs(t *testing.T) {
	type envVar struct {
		key   string
		value string
	}

	type GetUserNoEnvTest struct {
		Envs        []envVar
		GetUserTest GetUserTest
	}

	cases := []GetUserNoEnvTest{
		{
			Envs: []envVar{
				{
					key:   "MONGODB_ATLAS_HOST",
					value: "",
				},
				{
					key:   "MONGODB_ATLAS_USER",
					value: "user",
				},
				{
					key:   "MONGODB_ATLAS_PASS",
					value: "pass",
				},
			},
			GetUserTest: GetUserTest{
				verb:               "",
				path:               "",
				body:               nil,
				ExpectedStatusCode: 0,
				ExpectedErrorMsg:   "empty HOST, USER or PASS env vars for mongodb atlas",
			},
		},
		{
			Envs: []envVar{
				{
					key:   "MONGODB_ATLAS_HOST",
					value: "mongodb+srv://budget-tracker.gj4ww.mongodb.net",
				},
				{
					key:   "MONGODB_ATLAS_USER",
					value: "",
				},
				{
					key:   "MONGODB_ATLAS_PASS",
					value: "pass",
				},
			},
			GetUserTest: GetUserTest{
				verb:               "",
				path:               "",
				body:               nil,
				ExpectedStatusCode: 0,
				ExpectedErrorMsg:   "empty HOST, USER or PASS env vars for mongodb atlas",
			},
		},
		{
			Envs: []envVar{
				{
					key:   "MONGODB_ATLAS_HOST",
					value: "mongodb+srv://budget-tracker.gj4ww.mongodb.net",
				},
				{
					key:   "MONGODB_ATLAS_USER",
					value: "user",
				},
				{
					key:   "MONGODB_ATLAS_PASS",
					value: "",
				},
			},
			GetUserTest: GetUserTest{
				verb:               "",
				path:               "",
				body:               nil,
				ExpectedStatusCode: 0,
				ExpectedErrorMsg:   "empty HOST, USER or PASS env vars for mongodb atlas",
			},
		},
	}

	for _, testCase := range cases {
		for _, envVar := range testCase.Envs {
			t.Setenv(envVar.key, envVar.value)
		}

		_, err := NewClient()
		assert.Error(t, err, testCase.GetUserTest.ExpectedErrorMsg)
	}

}
