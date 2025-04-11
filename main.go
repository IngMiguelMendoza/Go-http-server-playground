package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const FILE_DIR_PATH = "./app"
	const SERVER_PORT = "8080"

	// Metrics handling
	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
	}

	// File system handling
	mux := http.NewServeMux()
	// mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(FILE_DIR_PATH)))))
	mux.Handle("/app/", apiCfg.metricsIncrement(http.StripPrefix("/app/", http.FileServer(http.Dir(FILE_DIR_PATH)))))

	// Health handling
	mux.HandleFunc("/healthz", readinessHandler)

	// Metrics handling
	mux.HandleFunc("/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/reset", apiCfg.resetMetricsHandler)

	server := &http.Server{
		Addr:    ":" + SERVER_PORT,
		Handler: mux,
	}

	log.Printf("HTTP Server starting at port: %s on filepath: %s", SERVER_PORT, FILE_DIR_PATH)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server listen and server err: %v", err)
	}

	log.Fatal("HTTP Server stopped")
}
