package auth

import (
	"net/http"
	"errors"
	"strings"
)

// define errors for reuse
var ErrNoAuthHeader = errors.New("authorization header missing")
var ErrInvalidAuthHeader = errors.New("invalid authorization header")

func GetBearerToken(headers http.Header) (string, error) {


	// lookup auth header
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}

	// trim whitespace and check prefix
	authHeader = strings.TrimSpace(authHeader)
	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return "", ErrInvalidAuthHeader
	}

	// strip bearer and return token
	token := strings.TrimSpace(authHeader[len("Bearer "):])
	if token == "" {
		return "", ErrInvalidAuthHeader      
	}

	return token, nil
}