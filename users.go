package main

import (
	"encoding/json"
	"net/http"
	"time"
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
	}
	
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
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