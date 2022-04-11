package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
)

func TestLoginHandler(t *testing.T) {
	type Query struct {
		urlPath    string
		input      objects.LoginMessage
		output     objects.TokenMessage
		statusCode int
	}

	tests := []struct {
		testName string
		queries  []Query
	}{
		{
			testName: "login_empty",
			queries: []Query{
				{
					urlPath:    "/api/user/login",
					input:      objects.LoginMessage{},
					statusCode: 400,
				},
				{
					urlPath:    "/api/user/register",
					input:      objects.LoginMessage{Login: "login"},
					statusCode: 400,
				},
				{
					urlPath:    "/api/user/register",
					input:      objects.LoginMessage{Password: "qwerty"},
					statusCode: 400,
				},
			},
		},
		{
			testName: "wrong_login",
			queries: []Query{
				{
					urlPath:    "/api/user/login",
					input:      objects.LoginMessage{Login: "exist_user", Password: "qwerty"},
					statusCode: 401,
				},
				{
					urlPath:    "/api/user/login",
					input:      objects.LoginMessage{Login: "exist_user2", Password: "correctpassword"},
					statusCode: 401,
				},
			},
		},
		{
			testName: "login_ok",
			queries: []Query{
				{
					urlPath:    "/api/user/login",
					input:      objects.LoginMessage{Login: "exist_user", Password: "correctpassword"},
					output:     objects.TokenMessage{Token: "token"},
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
				resp := testLoginRequest(t, ts, "POST", tq.urlPath, tq.input)
				defer resp.Body.Close()

				assert.Equal(t, tq.statusCode, resp.StatusCode)
				if resp.StatusCode == http.StatusOK {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						assert.Equal(t, err, nil)
						return
					}

					token := &objects.TokenMessage{}
					err = json.Unmarshal(body, token)
					if err != nil {
						assert.Equal(t, err, nil)
						return
					}
					assert.Equal(t, token.Token, tq.output.Token)
				}

				assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			}
		})
		ts.Close()
	}
}

func testLoginRequest(t *testing.T, ts *httptest.Server, method string, path string, input objects.LoginMessage) *http.Response {
	res, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewReader(res))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}
