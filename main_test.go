package main

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
)

func TestRouteStatusOK(t *testing.T) {

	routes := []string{
		"/",
		"/ping",
		"/v1/countries",
		"/v1/authors",
		"/v1/books",
		"/v1/customers",
		"/v1/publishers",
		"/v1/shipping-methods",
	}

	r := initRouter()

	for _, route := range routes {
		t.Run(route, func(t *testing.T) {
			req, _ := http.NewRequest("GET", route, nil)
			resp, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})

	}
}

func TestRouteStatusNotFound(t *testing.T) {
	routes := []string{
		"/foo",
		"/v1",
		"/v1/foo",
		"/v1/search/authors?name=foo",
		"/v1/search/books?author=foo",
		"/v1/search/books?title=foo",
		"/v1/search/books?isbn=5",
		"/v1/search/customer?email=foo",
	}

	r := initRouter()

	for _, route := range routes {
		t.Run(route, func(t *testing.T) {
			req, _ := http.NewRequest("GET", route, nil)
			resp, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		})
	}
}

func TestSearchErrors(t *testing.T) {
	var tests = []struct {
		search             string
		expectedStatusCode int
		expectedMessage    string
	}{
		{search: "/v1/customers/search?foo=1&bar=2", expectedStatusCode: fiber.ErrBadRequest.Code, expectedMessage: "multiple search terms not supported"},
		{search: "/v1/customers/search?foo=1", expectedStatusCode: fiber.ErrBadRequest.Code, expectedMessage: "no valid search term / value found. valid search terms: [email]"},
		{search: "/v1/books/search?foo=1&bar=2", expectedStatusCode: fiber.ErrBadRequest.Code, expectedMessage: "multiple search terms not supported"},
		{search: "/v1/books/search?foo=1", expectedStatusCode: fiber.ErrBadRequest.Code, expectedMessage: "no valid search term / value found. valid search terms: [title isbn author]"},
		{search: "/v1/authors/search?foo=1&bar=2", expectedStatusCode: fiber.ErrBadRequest.Code, expectedMessage: "multiple search terms not supported"},
		{search: "/v1/authors/search?foo=1", expectedStatusCode: fiber.ErrBadRequest.Code, expectedMessage: "no valid search term / value found. valid search terms: [name]"},
	}

	r := initRouter()

	for _, test := range tests {
		t.Run(test.search, func(t *testing.T) {
			req, _ := http.NewRequest("GET", test.search, nil)
			resp, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			var gr GravityResponse
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			json.Unmarshal(body, &gr)

			assert.Equal(t, test.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, test.expectedMessage, gr.Errors[0].Detail)
		})
	}
}

func TestLimit(t *testing.T) {

	var tests = []struct {
		route        string
		expectedSize int
	}{
		{route: "/v1/books", expectedSize: 100}, // current default maximum set to 100 in main.go
		{route: "/v1/books?limit=5", expectedSize: 5},
		{route: "/v1/books?limit=150", expectedSize: 100}, // checks that default maximum overrides excessive limit params
		{route: "/v1/books?limit=", expectedSize: 100},    // checks that default is applied if no value is provided
		{route: "/v1/books/search?author=Agatha Christie&limit=2", expectedSize: 2},
	}

	r := initRouter()

	for _, test := range tests {
		t.Run(test.route, func(t *testing.T) {
			req, _ := http.NewRequest("GET", test.route, nil)
			resp, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			a, _ := objx.FromJSON(string(body))

			assert.Len(t, a["data"], test.expectedSize)
		})
	}

}

func TestOffset(t *testing.T) {

	var tests = []struct {
		route           string
		expectedFirstId int
	}{
		{route: "/v1/books", expectedFirstId: 1},
		{route: "/v1/books?offset=5", expectedFirstId: 6},
		{route: "/v1/books?offset=150", expectedFirstId: 151},
		{route: "/v1/books?offset=", expectedFirstId: 1},
		{route: "/v1/books?offset=0", expectedFirstId: 1},
		{route: "/v1/books/search?author=Agatha Christie&offset=5", expectedFirstId: 9559},
	}

	r := initRouter()

	for _, test := range tests {
		t.Run(test.route, func(t *testing.T) {
			req, _ := http.NewRequest("GET", test.route, nil)
			resp, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			a, _ := objx.FromJSON(string(body))

			assert.Equal(t, test.expectedFirstId, a.Get("data[0].id").Int())
		})
	}

}

func TestCombinedLimitOffset(t *testing.T) {
	var tests = []struct {
		route           string
		expectedSize    int
		expectedFirstId int
	}{
		{route: "/v1/books?limit=50&offset=5", expectedSize: 50, expectedFirstId: 6},
		{route: "/v1/books?offset=5&limit=50", expectedSize: 50, expectedFirstId: 6},
		{route: "/v1/books/search?author=Agatha Christie&limit=10&offset=3", expectedSize: 10, expectedFirstId: 35},
	}

	r := initRouter()

	for _, test := range tests {
		t.Run(test.route, func(t *testing.T) {
			req, _ := http.NewRequest("GET", test.route, nil)
			resp, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			a, _ := objx.FromJSON(string(body))

			assert.Len(t, a["data"], test.expectedSize)
			assert.Equal(t, test.expectedFirstId, a.Get("data[0].id").Int())
		})
	}
}
