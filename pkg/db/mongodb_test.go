package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ func(int) *MongoDB = NewMongoDB

func TestMongoDBDisconnectHandlesZeroValue(t *testing.T) {
	mongodb := &MongoDB{}

	assert.NotPanics(t, mongodb.Disconnect)
}
