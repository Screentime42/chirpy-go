package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:		"chirpy-access",
		IssuedAt:	jwt.NewNumericDate(now),
		ExpiresAt: 	jwt.NewNumericDate(now.Add(expiresIn)),
		Subject: 	userID.String(),
	})

	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return ss, nil
}


