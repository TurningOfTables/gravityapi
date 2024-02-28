package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllCountriesError(t *testing.T) {
	db := connectToDb()
	c := initRouter().AcquireCtx()
	db.Close(context.Background())

	res, err := AllCountries(db, c)

	assert.Equal(t, []Country(nil), res)
	assert.Equal(t, "conn closed", err.Error())
}
