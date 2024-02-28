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

func TestAuthorSearch(t *testing.T) {
	var tests = []struct {
		search         string
		expectedAuthor string
	}{
		{search: "name=Agatha Christie", expectedAuthor: "Agatha Christie"},
	}

	r := initRouter()

	for _, test := range tests {
		t.Run(test.search, func(t *testing.T) {
			res, _ := http.NewRequest("GET", "/v1/authors/search?"+test.search, nil)
			resp, err := r.Test(res)
			if err != nil {
				t.Error(err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			a, _ := objx.FromJSON(string(body))

			assert.Equal(t, test.expectedAuthor, a.Get("data[0].authorName").Str())
		})
	}
}

func TestAllAuthorsError(t *testing.T) {
	db := connectToDb()
	c := initRouter().AcquireCtx()
	db.Close(context.Background())

	res, err := AllAuthors(db, c)

	assert.Equal(t, []Author(nil), res)
	assert.Equal(t, "conn closed", err.Error())
}

func TestAuthorSearchError(t *testing.T) {
	db := connectToDb()
	c := initRouter().AcquireCtx()
	db.Close(context.Background())

	res, err := AuthorsBySearchTerm(db, c, "name", "Agatha Christie")

	assert.Equal(t, []Author(nil), res)
	assert.Equal(t, "conn closed", err.Error())
}

func TestAuthorSearchInvalidTerm(t *testing.T) {
	var db *pgx.Conn
	var c fiber.Ctx

	res, err := AuthorsBySearchTerm(db, c, "foo", "bar")

	assert.Equal(t, []Author(nil), res)
	assert.Equal(t, errors.New("invalid search term"), err)
}
