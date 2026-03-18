package auth_test

import (
    "testing"

    "github.com/Screentime42/chirpy-go/internal/auth"
)

func TestHashAndCheckPassword(t *testing.T) {
    password := "supersecret123"

    hash, err := auth.HashPassword(password)
    if err != nil {
        t.Fatalf("HashPassword returned error: %v", err)
    }

    if hash == "" {
        t.Fatalf("expected non-empty hash")
    }

    match, err := auth.CheckPasswordHash(password, hash)
    if err != nil {
        t.Fatalf("CheckPasswordHash returned error: %v", err)
    }

    if !match {
        t.Fatalf("expected password to match hash")
    }
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
    hash, _ := auth.HashPassword("correct-password")

    match, err := auth.CheckPasswordHash("wrong-password", hash)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if match {
        t.Fatalf("expected mismatch but got match")
	 }
	}
    