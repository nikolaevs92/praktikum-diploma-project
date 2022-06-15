package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

func TestRegistraionHandler(t *testing.T) {
	type Query struct {
		urlPath    string
		input      objects.RegisterMessage
		statusCode int
	}

	tests := []struct {
		testName string
		queries  []Query
	}{
		{
			testName: "registration_empty",
			queries: []Query{
				{
					urlPath:    "/api/user/register",
					input:      objects.RegisterMessage{},
					statusCode: 400,
				},
				{
					urlPath:    "/api/user/register",
					input:      objects.RegisterMessage{Login: "login"},
					statusCode: 400,
				},
				{
					urlPath:    "/api/user/register",
					input:      objects.RegisterMessage{Password: "qwerty"},
					statusCode: 400,
				},
			},
		},
		{
			testName: "registration_exist",
			queries: []Query{
				{
					urlPath:    "/api/user/register",
					input:      objects.RegisterMessage{Login: "exist_user", Password: "qwerty"},
					statusCode: 409,
				},
			},
		},
		{
			testName: "registration_ok",
			queries: []Query{
				{
					urlPath:    "/api/user/register",
					input:      objects.RegisterMessage{Login: "notexist_user", Password: "qwerty"},
					statusCode: 200,
				},
				{
					urlPath:    "/api/user/register",
					input:      objects.RegisterMessage{Login: "r", Password: "qwerty3"},
					statusCode: 200,
				},
			},
		},
	}

	for _, tt := range tests {
		g := GofemartTest{}
		a := AuthorizationTest{}

		r := MakeRouter(g, a)
		ts := httptest.NewServer(r)

		t.Run(tt.testName, func(t *testing.T) {
			for _, tq := range tt.queries {
				resp := testRegistrationsRequest(t, ts, "POST", tq.urlPath, tq.input)
				defer resp.Body.Close()

				assert.Equal(t, tq.statusCode, resp.StatusCode)
				assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			}
		})
		ts.Close()
	}
}

func testRegistrationsRequest(t *testing.T, ts *httptest.Server, method string, path string, input objects.RegisterMessage) *http.Response {
	res, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewReader(res))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	_, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp
}
