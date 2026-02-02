package main

import (
	"encoding/json"
	"net/http"
	"strings"
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
func handlerValidate (w http.ResponseWriter, r *http.Request) {
	
	// Expected JSON payload structure.
	type parameters struct {
		Body string `json:"body"`
	}

	// Decode the incoming JSON body into params
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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

	// Respond with a simple JSON payload indicating success
	respondWithJSON(w, http.StatusOK, map[string]string{
		"cleaned_body": cleaned,
	})
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



