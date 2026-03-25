package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Screentime42/chirpy-go/internal/auth"
	"github.com/Screentime42/chirpy-go/internal/database"
	"github.com/google/uuid"
)

type User struct { 
	ID uuid.UUID `json:"id"` 
	CreatedAt time.Time `json:"created_at"` 
	UpdatedAt time.Time `json:"updated_at"` 
	Email string `json:"email"`
}


func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	
	var params struct {
		Email	string	`json:"email"`
		Password string `json:"password"`
	}

	
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:				params.Email,
		HashedPassword: 	hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	user := User{ 
		ID: 			dbUser.ID, 
		CreatedAt: 	dbUser.CreatedAt, 
		UpdatedAt: 	dbUser.UpdatedAt, 
		Email: 		dbUser.Email, 
	}

	respondWithJSON(w, http.StatusCreated, user)
}


func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {

	var params struct {
		Email					string	`json:"email"`
		Password 			string 	`json:"password"`
		ExpiresInSeconds	int		`json:"expires_in_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}


	dbUser, err := cfg.db.LookUpUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil || match == false {
		respondWithError(w,http.StatusUnauthorized, "Incorrect email or password")
		return
	}


	maxExpiry := 3600
	expirySeconds := params.ExpiresInSeconds

	if expirySeconds <= 0 {
		expirySeconds = maxExpiry
	} else {
		if expirySeconds > maxExpiry {
			expirySeconds = maxExpiry
		}
	}

	token, err := auth.MakeJWT(
		dbUser.ID,
		cfg.JWTSecret,
		time.Duration(expirySeconds) * time.Second,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token")
		return
	}

	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:		refreshToken,
		UserID:		dbUser.ID,
		ExpiresAt:	time.Now().Add(60 * 24 * time.Hour),
		})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create refresh token")
		return
	}
	
	type response struct {
		User
		Token				string	`json:"token"`
		RefreshToken 	string	`json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:	dbUser.ID,
			CreatedAt:	dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email: dbUser.Email,
		},
		Token: token, 
		RefreshToken: refreshToken,
	})
}


func (cfg *apiConfig) handlerUsersUpdate (w http.ResponseWriter, r *http.Request) {
	// Auth the user (extract token from header)
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find token")
		return
	}

	userID, err :l= cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}


	// Parse request body
	type request struct {
			Email		string	`json:"email"`
			Password string	`json:"password"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	
	// Password hash
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	// Update the user in the db
	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:					userID,
		Email:				req.Email,
		HashedPassword:	string(hashedPassword)
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not update user")
	}

}