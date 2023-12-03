package forecast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type response struct {
	Properties properties `json:"properties"`
}

type properties struct {
	Periods []period `json:"periods"`
}

type period struct {
	DetailedForecast string `json:"detailedForecast"`
	Name             string `json:"name"`
}

type errorResponse struct {
	Detail string `json:"detail"`
}

func GetForecastFromURL(url, prefix string, hourly bool) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	return processResponse(res, prefix, hourly)
}

func processResponse(res *http.Response, prefix string, hourly bool) (string, error) {
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

	sb := strings.Builder{}
	sb.Grow(2001)
	sb.WriteString(prefix)
	for _, period := range data.Properties.Periods {
		appendString := fmt.Sprintf("--%s: %s\n", period.Name, period.DetailedForecast)
		if sb.Len()+len(appendString) > 2001 {
			break
		}
		sb.WriteString(appendString)
	}
	result := sb.String()[:sb.Len()-1] // Remove last \n

	return result, nil
}
