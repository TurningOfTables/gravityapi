package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllPublishersError(t *testing.T) {
	db := connectToDb()
	c := initRouter().AcquireCtx()
	db.Close(context.Background())

	res, err := AllPublishers(db, c)

	assert.Equal(t, []Publisher(nil), res)
	assert.Equal(t, "conn closed", err.Error())
}
