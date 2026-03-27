package auth

import (
	"net/http"
	"errors"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", ErrNoAuthHeader
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", errors.New("incorrect Authorization format")
	}

	return parts[1], nil
}