package main

import (
	"chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	wjtSecret      string
}

func main() {
	const FILE_DIR_PATH = "./app"
	const SERVER_PORT = "8080"

	// SQL initialization/setup
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found, must be set")
	}

	platform := os.Getenv("PLATFORM")
	if dbURL == "" {
		log.Fatal("DB_URL not found, must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(db)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not found, must be set")
	}

	// Configuration handling
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		wjtSecret:      secret,
	}

	// File system handling
	mux := http.NewServeMux()
	mux.Handle("/app/",
		apiCfg.metricsIncrement(
			middlewareLog(
				http.StripPrefix("/app/",
					http.FileServer(http.Dir(FILE_DIR_PATH))))),
	)

	/*
		API Headers
	*/
	mux.HandleFunc("POST /api/users", apiCfg.usersHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.apiPostChirpHandler)
	mux.HandleFunc("POST /api/login", apiCfg.apiUserLoginHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.apiTokenRefreshHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.apiTokenRevokeHandler)

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.apiGetChirpHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.apiGetChirpByIdHandler)

	/*
		Administration
	*/

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetricsHandler)

	server := &http.Server{
		Addr:    ":" + SERVER_PORT,
		Handler: mux,
	}

	log.Printf("HTTP Server starting at port: %s on filepath: %s", SERVER_PORT, FILE_DIR_PATH)
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("Server listen and server err: %v", err)
	}

	log.Fatal("HTTP Server stopped")
}
