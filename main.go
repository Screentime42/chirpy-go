// Simple HTTP server that serves static files under /app,
// exposes health/metrics endpoints, and tracks file server hits.


package main

import (
	"log"
	"net/http"
	"sync/atomic"
	
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	
)




// apiConfig holds shared application state.
// fileserverHits tracks how many times the file server endpoint is accessed.
type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	godotenv.Load()

	// Application configuration constants.
	const filepathRoot = "."
	const port = "8080"

	// Initialize API configuration, including the atomic hit counter.
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	// Create a file server for /app and wrap it with middleware // that increments the hit counter on each request.
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	// Health check endpoint for readiness probes.
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Validate "Chirps" are within parameters
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	// Reset the metrics counter.
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// Expose current metrics (e.g., file server hit count).
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start the HTTP server.
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
