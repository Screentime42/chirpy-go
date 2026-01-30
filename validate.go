package main

import (
	"encoding/json"
	"net/http"
)

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

	// Respond with a simple JSON payload indicating success
	respondWithJSON(w, http.StatusOK, map[string]bool{
		"valid": true,
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



