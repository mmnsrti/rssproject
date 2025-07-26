package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("API key not provided")
	}
	vals := strings.Split(apiKey, " ")
	if len(vals) != 2 {
		return "", errors.New("invalid API key format, malformed header")

	}
	if vals[0] != "ApiKey" {
		return "", errors.New("invalid API key format, malformed header, expected 'ApiKey <key>'")
	}
	return vals[1], nil
}
