package forecast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type response struct {
	Properties struct {
		Periods []struct {
			DetailedForecast string `json:"detailedForecast"`
			Name             string `json:"name"`
		} `json:"periods"`
	} `json:"properties"`
}

func GetForecastFromURL(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode > 299 {
		return "", fmt.Errorf("error getting forecast url from forecast endpoint with status code: %d", res.StatusCode)
	}
	if err != nil {
		return "", err
	}

	data := &response{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return "", err
	}

	result := ""
	for _, period := range data.Properties.Periods {
		result += fmt.Sprintf("%s: %s\n", period.Name, period.DetailedForecast)
	}
	result = result[:len(result)-1]

	return result, nil
}
