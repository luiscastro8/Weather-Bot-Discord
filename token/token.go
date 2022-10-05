package token

import (
	"errors"
	"os"
)

func GetToken() (string, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return "", errors.New("DISCORD_TOKEN environment variable is empty")
	}
	return token, nil
}
