package discord

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZip(t *testing.T) {
	t.Run("it should return false with 3 digits", func(t *testing.T) {
		assert.False(t, isValidZip("378"))
	})
	t.Run("it should return false with 6 digits", func(t *testing.T) {
		assert.False(t, isValidZip("902634"))
	})
	t.Run("it should return false with 4 digits and 1 character", func(t *testing.T) {
		assert.False(t, isValidZip("73c94"))
	})
	t.Run("it should return true with 5 digits", func(t *testing.T) {
		assert.True(t, isValidZip("93520"))
	})
}
