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
					key:   "MONGODB_HOST",
					value: "",
				},
				{
					key:   "MONGODB_USER",
					value: "user",
				},
				{
					key:   "MONGODB_PASS",
					value: "pass",
				},
			},
			GetUserTest: GetUserTest{
				verb:               "",
				path:               "",
				body:               nil,
				ExpectedStatusCode: 0,
				ExpectedErrorMsg:   "empty MONGODB_HOST, MONGODB_USER or MONGODB_PASS env vars for mongodb",
			},
		},
		{
			Envs: []envVar{
				{
					key:   "MONGODB_HOST",
					value: "mongodb+srv://budget-tracker.gj4ww.mongodb.net",
				},
				{
					key:   "MONGODB_USER",
					value: "",
				},
				{
					key:   "MONGODB_PASS",
					value: "pass",
				},
			},
			GetUserTest: GetUserTest{
				verb:               "",
				path:               "",
				body:               nil,
				ExpectedStatusCode: 0,
				ExpectedErrorMsg:   "empty MONGODB_HOST, MONGODB_USER or MONGODB_PASS env vars for mongodb",
			},
		},
		{
			Envs: []envVar{
				{
					key:   "MONGODB_HOST",
					value: "mongodb+srv://budget-tracker.gj4ww.mongodb.net",
				},
				{
					key:   "MONGODB_USER",
					value: "user",
				},
				{
					key:   "MONGODB_PASS",
					value: "",
				},
			},
			GetUserTest: GetUserTest{
				verb:               "",
				path:               "",
				body:               nil,
				ExpectedStatusCode: 0,
				ExpectedErrorMsg:   "empty MONGODB_HOST, MONGODB_USER or MONGODB_PASS env vars for mongodb",
			},
		},
	}

	for _, testCase := range cases {
		for _, envVar := range testCase.Envs {
			t.Log(envVar)

			t.Setenv(envVar.key, envVar.value)
		}

		_, err := NewClient()
		t.Log(err)
		// In case of expected error msg, validate `err`
		if testCase.GetUserTest.ExpectedErrorMsg != "" {
			assert.EqualError(t, err, testCase.GetUserTest.ExpectedErrorMsg)
		} else {
			assert.NoError(t, err)
		}
	}
}
