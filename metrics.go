package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (config *apiConfig) metricsIncrement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		config.fileServerHits.Add(1)
		next.ServeHTTP(resp, req)
	})
}

func (config *apiConfig) metricsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(fmt.Sprintf("Hits: %v", config.fileServerHits.Load())))
}
