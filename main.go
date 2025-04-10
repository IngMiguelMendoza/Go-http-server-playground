package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("HTTP Server starting")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server listen and server err: %v", err)
	}

	fmt.Println("Server listening")
}
