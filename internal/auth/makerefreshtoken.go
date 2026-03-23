package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	token := make([]byte, 32)
	rand.Read(token)
	encodedToken := hex.EncodeToString(token)

	return encodedToken
}