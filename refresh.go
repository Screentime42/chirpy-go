package main

import (
	"net/http"
	"time"

	"github.com/Screentime42/chirpy-go/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find token")
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find token")
		return
	}


	newToken, err := auth.MakeJWT(user.UserID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token")
		return
	}

	type response struct {
		Token		string	`json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})
}