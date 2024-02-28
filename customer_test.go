package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
)

func TestCustomerSearch(t *testing.T) {
	var tests = []struct {
		search           string
		expectedCustomer string
	}{
		{search: "email=rvatini1@fema.gov", expectedCustomer: "Vatini"},
	}

	r := initRouter()

	for _, test := range tests {
		t.Run(test.search, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/customers/search?"+test.search, nil)
			resp, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			a, _ := objx.FromJSON(string(body))

			assert.Equal(t, test.expectedCustomer, a.Get("data[0].lastName").Str())
		})
	}
}

func TestAllCustomersError(t *testing.T) {
	db := connectToDb()
	c := initRouter().AcquireCtx()
	db.Close(context.Background())

	res, err := AllCustomers(db, c)

	assert.Equal(t, []Customer(nil), res)
	assert.Equal(t, "conn closed", err.Error())
}

func TestCustomerSearchError(t *testing.T) {
	db := connectToDb()
	c := initRouter().AcquireCtx()
	db.Close(context.Background())

	res, err := CustomersBySearchTerm(db, c, "email", "rvatini1@fema.gov")

	assert.Equal(t, []Customer(nil), res)
	assert.Equal(t, "conn closed", err.Error())
}

func TestCustomerSearchInvalidTerm(t *testing.T) {
	var db *pgx.Conn
	var c fiber.Ctx

	res, err := CustomersBySearchTerm(db, c, "foo", "bar")

	assert.Equal(t, []Customer(nil), res)
	assert.Equal(t, errors.New("invalid search term"), err)
}
