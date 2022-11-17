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

type errorResponse struct {
	Detail string `json:"detail"`
}

func GetForecastFromURL(url, prefix string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	return processResponse(res, prefix)
}

func processResponse(res *http.Response, prefix string) (string, error) {
	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode > 299 {
		data := &errorResponse{}
		err = json.Unmarshal(body, data)
		if err != nil {
			return "", fmt.Errorf("error getting forecast from forecast endpoint with status code %d. Unable to parse error detail", res.StatusCode)
		}
		return "", fmt.Errorf("error getting forecast from forecast endpoint with status code %d and reason: %s", res.StatusCode, data.Detail)
	}
	if err != nil {
		return "", err
	}

	data := &response{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return "", err
	}

	result := prefix
	for _, period := range data.Properties.Periods {
		appendString := fmt.Sprintf("--%s: %s\n", period.Name, period.DetailedForecast)
		if len(result)+len(appendString) > 2001 {
			break
		}
		result += appendString
	}
	result = result[:len(result)-1] // Remove last \n

	return result, nil
}
