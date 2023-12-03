package points

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type response struct {
	Properties properties `json:"properties"`
}

type properties struct {
	Forecast       string `json:"forecast"`
	ForecastHourly string `json:"forecastHourly"`
}

type errorResponse struct {
	Detail string `json:"detail"`
}

func GetDailyForecastURLFromCoords(lat, long string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%s,%s", lat, long))
	if err != nil {
		return "", err
	}

	return processResponse(res, false)
}

func GetHourlyForecastURLFromCoords(lat, long string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%s,%s", lat, long))
	if err != nil {
		return "", err
	}

	return processResponse(res, true)
}

func processResponse(res *http.Response, hourly bool) (string, error) {
	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode > 299 {
		data := &errorResponse{}
		err = json.Unmarshal(body, data)
		if err != nil {
			return "", fmt.Errorf("error getting forecast url from points endpoint with status code %d. Unable to parse error detail", res.StatusCode)
		}
		return "", fmt.Errorf("error getting forecast url from points endpoint with status code %d and reason: %s", res.StatusCode, data.Detail)
	}
	if err != nil {
		return "", err
	}

	data := &response{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return "", err
	}

	if hourly {
		return data.Properties.ForecastHourly, nil
	}
	return data.Properties.Forecast, nil
}
