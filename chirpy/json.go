package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithErrorJson(resp http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Print(err)
	}

	if code > 499 {
		log.Printf("Responding with server error, 5xx: %s", msg)
	}

	type errorValue struct {
		Error string `json:"error"`
	}

	respondWithJSON(resp, code, errorValue{
		Error: msg,
	})
}

func respondWithJSON(resp http.ResponseWriter, code int, payload interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		resp.WriteHeader(500)
		return
	}

	resp.WriteHeader(code)
	resp.Write(data)
}
