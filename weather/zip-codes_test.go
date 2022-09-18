package weather

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testZipFile = "test-files/zip-codes-test.csv"

func TestGetCoordsFromZip(t *testing.T) {
	t.Run("Get coordinates from valid zip code", func(t *testing.T) {
		err := OpenZipFile(testZipFile)
		assert.Nil(t, err)

		latitude, longitude, err := GetCoordsFromZip("99571")
		assert.Nil(t, err)
		assert.Equal(t, "55.1858", latitude)
		assert.Equal(t, "-162.7211", longitude)
	})

	t.Run("Return error if zip code is not in database", func(t *testing.T) {
		err := OpenZipFile(testZipFile)
		assert.Nil(t, err)

		_, _, err = GetCoordsFromZip("83729")
		assert.Error(t, err)
	})

	t.Run("Return error if zip cache has not been initialized", func(t *testing.T) {
		zipCodes = nil

		latitude, longitude, err := GetCoordsFromZip("99571")
		assert.Error(t, err)
		assert.NotEqual(t, "55.1858", latitude)
		assert.NotEqual(t, "-162.7211", longitude)
	})
}

func TestOpenZipFile(t *testing.T) {
	t.Run("Open file without error", func(t *testing.T) {
		err := OpenZipFile(testZipFile)
		assert.Nil(t, err)
	})

	t.Run("Return error if file name is invalid", func(t *testing.T) {
		err := OpenZipFile("fakefile.banana")
		assert.Error(t, err)
	})
}
