package main

import (
	"log"
	"net/http"
)

func main() {
	const FILE_DIR_PATH = "."
	const SERVER_PORT = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(FILE_DIR_PATH)))

	server := &http.Server{
		Addr:    ":" + SERVER_PORT,
		Handler: mux,
	}

	log.Println("HTTP Server starting at port: %s", SERVER_PORT)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server listen and server err: %v", err)
	}

	log.Fatal("HTTP Server stopped")
}
