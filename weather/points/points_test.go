package points

import (
	"Weather-Bot-Discord/testutil"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestProcessResponse(t *testing.T) {
	t.Run("It should return an error if the status code was not 2XX", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(strings.NewReader("")),
		}
		s, err := processResponse(mockResponse)
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return an error if the status code was not 2XX and the response has error details", func(t *testing.T) {
		bodyBytes, err := json.Marshal(errorResponse{Detail: "an expected error has occurred"})
		if err != nil {
			panic(err)
		}
		mockResponse := &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(strings.NewReader(string(bodyBytes))),
		}
		s, err := processResponse(mockResponse)
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return an error if unable to read from response body", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(testutil.NewErrorReader()),
		}
		s, err := processResponse(mockResponse)
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return an error if unable to unmarshal json from response body", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("thisisnotavalidresponse")),
		}
		s, err := processResponse(mockResponse)
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return the forecast url from the response", func(t *testing.T) {
		bodyBytes, err := json.Marshal(response{Properties: properties{Forecast: "https://api.com/234"}})
		if err != nil {
			panic(err)
		}
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(string(bodyBytes))),
		}
		s, err := processResponse(mockResponse)
		assert.Equal(t, "https://api.com/234", s)
		assert.Nil(t, err)
	})
}
