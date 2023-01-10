package forecast

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
		mockResponse := createMockResponseWithString("", 500)
		s, err := processResponse(mockResponse, "")
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return an error if the status code was not 2XX and the response has error details", func(t *testing.T) {
		responseBody := createMockResponseErrorBody(errorResponse{Detail: "an expected error has occurred"})
		mockResponse := createMockResponseWithString(responseBody, 500)
		s, err := processResponse(mockResponse, "")
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return an error if unable to read from response body", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(iotest.ErrReader(errors.New("intentional error"))),
		}
		s, err := processResponse(mockResponse, "")
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return an error if unable to unmarshal json from response body", func(t *testing.T) {
		mockResponse := createMockResponseWithString("thisisnotavalidstruct", 200)
		s, err := processResponse(mockResponse, "")
		assert.Equal(t, "", s)
		assert.Error(t, err)
	})

	t.Run("It should return the weather with a prefix", func(t *testing.T) {
		responseBody := response{Properties: properties{Periods: []period{{
			DetailedForecast: "cloudy with a chance of water falling",
			Name:             "Sunday",
		}}}}
		mockResponse := createMockResponse(responseBody, 200)
		s, err := processResponse(mockResponse, "This is the weather for 1234 lane ln\n")
		assert.Equal(t, "This is the weather for 1234 lane ln\n--Sunday: cloudy with a chance of water falling", s)
		assert.Nil(t, err)
	})

	t.Run("It should return the weather for one day", func(t *testing.T) {
		responseBody := response{Properties: properties{Periods: []period{{
			DetailedForecast: "cloudy with a chance of rain",
			Name:             "Monday",
		}}}}
		mockResponse := createMockResponse(responseBody, 200)
		s, err := processResponse(mockResponse, "")
		assert.Equal(t, "--Monday: cloudy with a chance of rain", s)
		assert.Nil(t, err)
	})

	t.Run("It should return the weather for two days", func(t *testing.T) {
		responseBody := response{Properties: properties{Periods: []period{{
			DetailedForecast: "snowy",
			Name:             "Tuesday",
		}, {
			DetailedForecast: "sunny and hot",
			Name:             "Wednesday",
		}}}}
		mockResponse := createMockResponse(responseBody, 200)
		s, err := processResponse(mockResponse, "")
		assert.Equal(t, "--Tuesday: snowy\n--Wednesday: sunny and hot", s)
		assert.Nil(t, err)
	})

	t.Run("It should return the weather without exceeding 2000 characters", func(t *testing.T) {
		longForecast := make([]byte, 2005)
		for i := range longForecast {
			longForecast[i] = 't'
		}
		responseBody := response{Properties: properties{Periods: []period{{
			DetailedForecast: "volcanic eruption",
			Name:             "Thursday",
		}, {
			DetailedForecast: string(longForecast),
			Name:             "Friday",
		}}}}
		mockResponse := createMockResponse(responseBody, 200)
		s, err := processResponse(mockResponse, "")
		assert.Equal(t, "--Thursday: volcanic eruption", s)
		assert.Nil(t, err)
	})
}

func createMockResponseWithString(body string, code int) *http.Response {
	return &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func createMockResponse(body response, code int) *http.Response {
	bytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return createMockResponseWithString(string(bytes), code)
}

func createMockResponseErrorBody(r errorResponse) string {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
