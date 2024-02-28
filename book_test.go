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

func TestBookSearch(t *testing.T) {
	var tests = []struct {
		search        string
		expectedTitle string
	}{
		{search: "title=The Tempest", expectedTitle: "The Tempest"},
		{search: "isbn=9781559277587", expectedTitle: "They Do It With Mirrors"},
		{search: "author=Agatha Christie", expectedTitle: "Hercule Poirot's Christmas: A BBC Radio 4 Full-Cast Dramatisation"},
	}

	r := initRouter()

	for _, test := range tests {
		t.Run(test.search, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/books/search?"+test.search, nil)
			resp, err := r.Test(req)
			if err != nil {
				t.Error(err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			a, _ := objx.FromJSON(string(body))

			assert.Equal(t, test.expectedTitle, a.Get("data[0].title").Str())
		})
	}

}

func TestAllBooksError(t *testing.T) {
	db := connectToDb()
	c := initRouter().AcquireCtx()
	db.Close(context.Background())

	res, err := AllBooks(db, c)

	assert.Equal(t, []Book(nil), res)
	assert.Equal(t, "conn closed", err.Error())
}

func TestBookSearchError(t *testing.T) {
	db := connectToDb()
	c := initRouter().AcquireCtx()
	db.Close(context.Background())

	res, err := BooksBySearchTerm(db, c, "title", "The Tempest")

	assert.Equal(t, []Book(nil), res)
	assert.Equal(t, "conn closed", err.Error())
}

func TestBookSearchInvalidTerm(t *testing.T) {
	var db *pgx.Conn
	var c fiber.Ctx

	res, err := BooksBySearchTerm(db, c, "foo", "bar")

	assert.Equal(t, []Book(nil), res)
	assert.Equal(t, errors.New("invalid search term"), err)
}
