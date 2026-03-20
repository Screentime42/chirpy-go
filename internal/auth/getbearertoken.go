package auth

import (
	"net/http"
	"errors"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	// define errors for reuse
	ErrNoAuthHeader := errors.New("authorization header missing")
	ErrInvalidAuthHeader := errors.New("invalid authorization header")

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