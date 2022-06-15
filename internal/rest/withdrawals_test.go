package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nikolaevs92/praktikum-diploma-project.git/internal/objects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostWithdrawHandler(t *testing.T) {
	tests := []struct {
		testName   string
		input      objects.Withdraw
		statusCode int
	}{
		{
			testName: "test200",
			input: objects.Withdraw{
				Order: "134",
				Sum:   5,
			},
			statusCode: 200,
		},
		{
			testName: "test402",
			input: objects.Withdraw{
				Order: "134",
				Sum:   105,
			},
			statusCode: 402,
		},
		{
			testName: "test422",
			input: objects.Withdraw{
				Order: "wrong",
				Sum:   5,
			},
			statusCode: 422,
		},
	}

	g := GofemartTest{}
	a := AuthorizationTest{}

	r := MakeRouter(g, a)
	ts := httptest.NewServer(r)
	for _, tt := range tests {

		t.Run(tt.testName, func(t *testing.T) {
			resp := testPostWithdrawRequest(t, ts, "POST", "/api/user/balance/withdraw", tt.input)
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
	ts.Close()
}

func testPostWithdrawRequest(t *testing.T, ts *httptest.Server, method string, path string, input objects.Withdraw) *http.Response {
	res, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewReader([]byte(res)))
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+"token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

func TestGetWitdrawalshHandler(t *testing.T) {
	g := GofemartTest{}
	a := AuthorizationTest{}

	r := MakeRouter(g, a)
	ts := httptest.NewServer(r)
	t.Run("test not empty withdrawals", func(t *testing.T) {
		resp := testGetWithdrawalsRequest(t, ts, "GET", "/api/user/balance/withdrawals", "token")
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		withdrawals := []objects.Withdraw{}
		err = json.Unmarshal(body, &withdrawals)
		assert.NoError(t, err)

		assert.Equal(t, len(withdrawals), 1)
		assert.Equal(t, withdrawals[0].Sum, 40.0)
		assert.Equal(t, withdrawals[0].Order, "134")
	})

	t.Run("test empty withdrawals", func(t *testing.T) {
		resp := testGetWithdrawalsRequest(t, ts, "GET", "/api/user/balance/withdrawals", "token3")
		defer resp.Body.Close()

		assert.Equal(t, 204, resp.StatusCode)
	})

	ts.Close()
}

func testGetWithdrawalsRequest(t *testing.T, ts *httptest.Server, method string, path string, token string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewReader([]byte("")))
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}
