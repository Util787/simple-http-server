package auth

import (
	"errors"
	"net/http"
)

// expecting header to be in this format: "Authorization:ApiKey THE_KEY_HERE"
func GetAPIKey(headers http.Header) (string, error) {
	APIKeyString := headers.Get("Authorization")
	if APIKeyString == "" {
		return "", errors.New("empty Authorization header") // failed to get token string:
	}
	return APIKeyString[7:], nil
}
