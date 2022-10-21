package zip

import (
	"encoding/csv"
	"errors"
	"os"
)

var zipCodes map[string][]string

func GetCoordsFromZip(zip string) (string, string, error) {
	if zipCodes == nil {
		return "", "", errors.New("zip code cache is empty")
	}

	coords, present := zipCodes[zip]
	if !present {
		return "", "", errors.New("zip code is not in cache")
	}
	return coords[0], coords[1], nil
}

func OpenZipFile(fileName string) error {
	open, err := os.Open(fileName)
	if err != nil {
		return err
	}

	reader := csv.NewReader(open)
	data, err := reader.ReadAll()
	if err != nil {
		return err
	}

	zipCodes = make(map[string][]string, len(data))
	for _, val := range data {
		zipCodes[val[0]] = []string{val[1], val[2]}
	}
	return nil
}
