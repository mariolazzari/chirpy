package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Returns API key retrieved from Authorization header
func GetAPIKey(header http.Header) (string, error) {
	apiKey := header.Get("Authorization")

	if apiKey == "" {
		return "", errors.New("API Key is required but got none")
	}

	return strings.TrimPrefix(apiKey, "ApiKey "), nil
}
