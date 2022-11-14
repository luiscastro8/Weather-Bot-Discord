package myerrors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAddressNotFoundError(t *testing.T) {
	t.Run("It should return error with correct message and unmatched address", func(t *testing.T) {
		testMessage := "this is a test message"
		address := "12345678 fake street ln"
		err := NewAddressNotFoundError(testMessage, address)
		assert.Equal(t, testMessage, err.Error())
		assert.Equal(t, address, err.UnmatchedAddress)
	})
}
