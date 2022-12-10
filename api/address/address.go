package address

import (
	"Weather-Bot-Discord/myerrors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type response struct {
	Result struct {
		AddressMatches []struct {
			Coordinates struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"coordinates"`
			MatchedAddress string `json:"matchedAddress"`
		} `json:"addressMatches"`
	} `json:"result"`
}

func GetCoordsFromAddress(address string) (string, string, string, error) {
	res, err := http.Get(fmt.Sprintf("https://geocoding.geo.census.gov/geocoder/locations/onelineaddress?address=%s&benchmark=2020&format=json", url.QueryEscape(address)))
	if err != nil {
		return "", "", "", err
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode > 299 {
		return "", "", "", fmt.Errorf("error getting coordinates from geocoding endpoint with status code %d", res.StatusCode)
	}
	if err != nil {
		return "", "", "", err
	}

	data := &response{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return "", "", "", err
	}

	if len(data.Result.AddressMatches) == 0 {
		return "", "", "", myerrors.NewAddressNotFoundError("address not found for input "+address, address)
	}

	matchedAddress := data.Result.AddressMatches[0]
	return fmt.Sprintf("%.4f", matchedAddress.Coordinates.Y), fmt.Sprintf("%.4f", matchedAddress.Coordinates.X), matchedAddress.MatchedAddress, nil
}
