package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	e := &GravityError{Detail: "foo"}
	res := e.Error()
	assert.Equal(t, "foo", res)
}
