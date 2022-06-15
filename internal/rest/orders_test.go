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

func TestPostOrdersHandler(t *testing.T) {
	tests := []struct {
		testName   string
		input      string
		statusCode int
	}{
		{
			testName:   "test200",
			input:      "order200",
			statusCode: 200,
		},
		{
			testName:   "test202",
			input:      "order202",
			statusCode: 202,
		},
		{
			testName:   "test400",
			input:      "order400",
			statusCode: 400,
		},
		{
			testName:   "test401",
			input:      "order401",
			statusCode: 401,
		},
		{
			testName:   "test409",
			input:      "order409",
			statusCode: 409,
		},
		{
			testName:   "test422",
			input:      "order422",
			statusCode: 422,
		},
		{
			testName:   "test500",
			input:      "order500",
			statusCode: 500,
		},
	}

	g := GofemartTest{}
	a := AuthorizationTest{}

	r := MakeRouter(g, a)
	ts := httptest.NewServer(r)
	for _, tt := range tests {

		t.Run(tt.testName, func(t *testing.T) {
			resp := testPostOrdersRequest(t, ts, "POST", "/api/user/orders", tt.input)
			defer resp.Body.Close()

			assert.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
	ts.Close()
}

func testPostOrdersRequest(t *testing.T, ts *httptest.Server, method string, path string, input string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewReader([]byte(input)))
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+"token")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

func TestGetOrdersHandler(t *testing.T) {
	g := GofemartTest{}
	a := AuthorizationTest{}

	r := MakeRouter(g, a)
	ts := httptest.NewServer(r)
	t.Run("test not empty orders", func(t *testing.T) {
		resp := testGettOrdersRequest(t, ts, "GET", "/api/user/orders", "token")
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		orders := []objects.Order{}
		err = json.Unmarshal(body, &orders)
		assert.NoError(t, err)

		assert.Equal(t, len(orders), 1)
		assert.Equal(t, orders[0].Accural, 111.0)
		assert.Equal(t, orders[0].Status, objects.OrderStatusNew)
		assert.Equal(t, orders[0].Number, "124")
	})

	t.Run("test  empty orders", func(t *testing.T) {
		resp := testGettOrdersRequest(t, ts, "GET", "/api/user/orders", "token3")
		defer resp.Body.Close()

		assert.Equal(t, 204, resp.StatusCode)
	})

	ts.Close()
}

func testGettOrdersRequest(t *testing.T, ts *httptest.Server, method string, path string, token string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewReader([]byte("")))
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}
