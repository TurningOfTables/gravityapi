package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteStatusOK(t *testing.T) {

	routes := []string{
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

			assert.Equal(t, 200, resp.StatusCode)
		})

	}
}
