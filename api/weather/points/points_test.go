package points

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"
)

func TestProcessResponse(t *testing.T) {
	t.Run("It should return an error if the status code was not 2XX", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(strings.NewReader("")),
		}
		s, err := processResponse(mockResponse, false)
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
		s, err := processResponse(mockResponse, false)
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return an error if unable to read from response body", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(iotest.ErrReader(errors.New("intentional error"))),
		}
		s, err := processResponse(mockResponse, false)
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return an error if unable to unmarshal json from response body", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("thisisnotavalidresponse")),
		}
		s, err := processResponse(mockResponse, false)
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return the daily forecast url from the response", func(t *testing.T) {
		bodyBytes, err := json.Marshal(response{Properties: properties{Forecast: "https://api.com/234/daily", ForecastHourly: "https://api.com/234/hourly"}})
		if err != nil {
			panic(err)
		}
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(string(bodyBytes))),
		}
		s, err := processResponse(mockResponse, false)
		assert.Equal(t, "https://api.com/234/daily", s)
		assert.Nil(t, err)
	})

	t.Run("It should return the hourly forecast url from the response", func(t *testing.T) {
		bodyBytes, err := json.Marshal(response{Properties: properties{Forecast: "https://api.com/234/daily", ForecastHourly: "https://api.com/234/hourly"}})
		if err != nil {
			panic(err)
		}
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(string(bodyBytes))),
		}
		s, err := processResponse(mockResponse, true)
		assert.Equal(t, "https://api.com/234/hourly", s)
		assert.Nil(t, err)
	})
}
