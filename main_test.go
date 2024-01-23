package main

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestRouteStatusOK(t *testing.T) {

	routes := []string{
		"/ping",
		"/v1/countries",
		"/v1/authors",
		"/v1/authors?name=Abigail Adams",
		"/v1/books",
		"/v1/books?author=William Shakespeare",
		"/v1/books?title=The Tempest",
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
