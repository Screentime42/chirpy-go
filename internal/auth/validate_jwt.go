package auth

import (
	"errors"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}


	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID in token")
	}
	
	return userID, nil
}