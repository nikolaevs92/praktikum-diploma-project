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

func TestAuthToken(t *testing.T) {
	g := GofemartTest{}
	a := AuthorizationTest{}

	r := MakeRouter(g, a)
	ts := httptest.NewServer(r)

	t.Run("get token", func(t *testing.T) {
		resp := testLoginRequest(t, ts, "POST", "/api/user/login", objects.LoginMessage{Login: "exist_user", Password: "correctpassword"})
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
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
			assert.Equal(t, token.Token, "token")
		}
	})

	t.Run("correct token using", func(t *testing.T) {
		resp := testRequest(t, ts, "POST", "/api/user/orders", "correct_order", "token")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("incorrect token using", func(t *testing.T) {
		resp := testRequest(t, ts, "POST", "/api/user/orders", "correct_order", "token2")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
	ts.Close()
}

func testRequest(t *testing.T, ts *httptest.Server, method string, path string, input string, token string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewReader([]byte(input)))
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}
