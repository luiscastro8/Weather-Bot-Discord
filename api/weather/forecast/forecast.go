package forecast

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type dailyForecastResponse struct {
	Properties struct {
		Periods []struct {
			DetailedForecast string `json:"detailedForecast"`
			Name             string `json:"name"`
		} `json:"periods"`
	} `json:"properties"`
}

type hourlyForecastResponse struct {
	Properties struct {
		Periods []struct {
			StartTime                  string `json:"startTime"`
			Temperature                int    `json:"temperature"`
			TemperatureUnit            string `json:"temperatureUnit"`
			ProbabilityOfPrecipitation struct {
				Value int `json:"value"`
			} `json:"probabilityOfPrecipitation"`
			ShortForecast string `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
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

	if hourly {
		return getHourlyForecastMessage(body, prefix)
	}
	return getDailyForecastMessage(body, prefix)
}

func getDailyForecastMessage(body []byte, prefix string) (string, error) {
	data := &dailyForecastResponse{}
	err := json.Unmarshal(body, data)
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

func getHourlyForecastMessage(body []byte, prefix string) (string, error) {
	data := &hourlyForecastResponse{}
	err := json.Unmarshal(body, data)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}
	sb.Grow(2001)
	sb.WriteString(prefix)
	for _, period := range data.Properties.Periods {
		startTime, err := time.Parse(time.RFC3339, period.StartTime)
		if err != nil {
			return "", err
		}
		appendString := "--" + startTime.Format("Mon, 02 Jan 2006 03:04:05 PM MST") + "\n"
		appendString += "Short Forecast: " + period.ShortForecast + "\n"
		appendString += "Temperature: " + strconv.Itoa(period.Temperature) + "°" + period.TemperatureUnit + "\n"
		appendString += "Precipitation: " + strconv.Itoa(period.ProbabilityOfPrecipitation.Value) + "%\n"
		if sb.Len()+len(appendString) > 2001 {
			break
		}
		sb.WriteString(appendString)
	}
	result := sb.String()[:sb.Len()-1] // Remove last \n

	return result, nil
}
