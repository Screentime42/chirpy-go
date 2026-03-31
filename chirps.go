package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sort"

	"github.com/Screentime42/chirpy-go/internal/auth"
	"github.com/Screentime42/chirpy-go/internal/database"
	"github.com/google/uuid"
)

// List to hold banned words
var bannedWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert": {},
	"fornax": {},
}


// function to replace banned words with specified replacement
func censorBannedWords(body string, banned map[string]struct{}, replacement string) string {
	words := strings.Fields(body)
	out := make([]string, len(words))

	for i, w := range words {
		lower := strings.ToLower(w)

		if _, banned := banned[lower]; banned {
			out[i] = replacement
		} else {
			out[i] = w
		}
	}
	return strings.Join(out, " ")
}


// handlerValidate validates a chirp sent in the request body
// It ensures the JSON is valid and the chirp meets length requirements
func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {

	// user validation
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Expected JSON payload structure.
	type parameters struct {
		Body 		string 		`json:"body"`
	}

	// Decode the incoming JSON body into params
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		// Return a generic error to the client if JSON decoding fails
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Enforce the 140‑character chirp limit
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Apply censor
	cleaned := censorBannedWords(params.Body, bannedWords, "****")

	// Insert chirp to DB
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: 	cleaned,
		UserID: 	userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}


// respondWithError sends a JSON‑formatted error message with the given status code 
// Used for consistent error responses across the API
func respondWithError(w http.ResponseWriter, code int, msg string) {
   w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
    
	type errorResponse struct {
		Error string `json:"error"`
	}

	resp := errorResponse{
		Error: msg,
	}

	json.NewEncoder(w).Encode(resp)
}

// respondWithJSON sends any value as a JSON response with the given status code.
// Used for consistent success responses across the API
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(payload)
}



func (cfg *apiConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// Get author id from url query
	authorID := r.URL.Query().Get("author_id")
	
	// If authorID is present execute otherwise skip
	if authorID != "" {
		// Parse authorID from string to uuid
		id, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "could not parse id")
			return
		}

		// Return chirps that match with the authorid
		chirps, err := cfg.db.GetChirpsByAuthorID(r.Context(), id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not get chirps")
			return
		}
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get chirps")
		return
	}

	// Get sort instruction
	sortValue := r.URL.Query().Get("sort")
	
	// If sortValue is present execute otherwise skip. As chirps are default ASC sorted, only sorted if sortValue is "desc"
	if sortValue != "" {
		if sortValue == "desc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
			})
		}
	}

	respondWithJSON(w, http.StatusOK, chirps)
}


func (cfg *apiConfig) HandlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil { 
		respondWithError(w, http.StatusBadRequest, "invalid chirp id") 
		return
	}

	chirp, err := cfg.db.GetSingleChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}


func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	// Extract chirp id
	chirpIDStr := r.PathValue("chirp_id")

	// Convert to UUID
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}


	// Auth the user (extract token from header)
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}


	// Fetch chirp for ownership check
	chirp, err := cfg.db.GetSingleChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}

	// Check ownership of chirp
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "you cannot delete this chirp")
		return
	}


	// Delete chirp
	err = cfg.db.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{
		ID:		chirpID,
		UserID:	userID,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not delete chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}