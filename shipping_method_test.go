package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllShippingMethodsError(t *testing.T) {
	db := connectToDb()
	c := initRouter().AcquireCtx()
	db.Close(context.Background())

	res, err := AllShippingMethods(db, c)

	assert.Equal(t, []ShippingMethod(nil), res)
	assert.Equal(t, "conn closed", err.Error())
}
