package auth_test

import (
    "testing"
    "time"

    "github.com/google/uuid"
    "github.com/Screentime42/chirpy-go/internal/auth"
)

func TestMakeAndValidateJWT(t *testing.T) {
    secret := "mysecretkey"
    userID := uuid.New()

    token, err := auth.MakeJWT(userID, secret, time.Hour)
    if err != nil {
        t.Fatalf("MakeJWT returned error: %v", err)
    }

    returnedID, err := auth.ValidateJWT(token, secret)
    if err != nil {
        t.Fatalf("ValidateJWT returned error: %v", err)
    }

    if returnedID != userID {
        t.Fatalf("expected userID %v, got %v", userID, returnedID)
    }
}

func TestValidateJWT_InvalidSignature(t *testing.T) {
    secret := "correct-secret"
    wrongSecret := "wrong-secret"
    userID := uuid.New()

    token, err := auth.MakeJWT(userID, secret, time.Hour)
    if err != nil {
        t.Fatalf("MakeJWT returned error: %v", err)
    }

    _, err = auth.ValidateJWT(token, wrongSecret)
    if err == nil {
        t.Fatalf("expected error for invalid signature but got none")
    }
}

func TestValidateJWT_Expired(t *testing.T) {
    secret := "mysecretkey"
    userID := uuid.New()

    // Token expires immediately
    token, err := auth.MakeJWT(userID, secret, -1*time.Second)
    if err != nil {
        t.Fatalf("MakeJWT returned error: %v", err)
    }

    _, err = auth.ValidateJWT(token, secret)
    if err == nil {
        t.Fatalf("expected expiration error but got none")
    }
}
