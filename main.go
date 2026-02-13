// Simple HTTP server that serves static files under /app,
// exposes health/metrics endpoints, and tracks file server hits.


package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"os"
	"database/sql"
	
	"github.com/joho/godotenv"
	"github.com/Screentime42/chirpy-go/internal/database"
	_ "github.com/lib/pq"
	
)




// apiConfig holds shared application state.
// fileserverHits tracks how many times the file server endpoint is accessed.
type apiConfig struct {
	fileserverHits atomic.Int32
	db					*database.Queries
	Platform			string
}

func main() {
	// Application configuration constants.
	const filepathRoot = "."
	const port = "8080"

	// Loading database
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = "prod"
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: 	atomic.Int32{},
		db:             	dbQueries,
		Platform:			platform,
	}


	mux := http.NewServeMux()

	// Create a file server for /app and wrap it with middleware 
	// that increments the hit counter on each request.
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	// Health check endpoint for readiness probes.
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Reset the metrics counter.
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// Expose current metrics (e.g., file server hit count).
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	// Create new user handler link
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	// Validate and create new chirp
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Start the HTTP server.
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
