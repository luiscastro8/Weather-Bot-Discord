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
	Forecast string `json:"forecast"`
}

type errorResponse struct {
	Detail string `json:"detail"`
}

func GetForecastURLFromCoords(lat, long string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://api.weather.gov/points/%s,%s", lat, long))
	if err != nil {
		return "", err
	}

	return processResponse(res)
}

func processResponse(res *http.Response) (string, error) {
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

	return data.Properties.Forecast, nil
}
