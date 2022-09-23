package points

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type response struct {
	Properties struct {
		Forecast string `json:"forecast"`
	} `json:"properties"`
}

func GetForecastURLFromCoords(lat, long string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%s,%s", lat, long))
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode > 299 {
		return "", errors.New(fmt.Sprintf("error getting forecast url from points endpoint with status code: %d", res.StatusCode))
	}
	if err != nil {
		return "", err
	}

	data := &response{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return "", err
	}

	return data.Properties.Forecast, nil
}
