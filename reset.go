package main

import "net/http"

func (config *apiConfig) resetMetricsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	config.fileServerHits.Store(0)
	resp.Write([]byte(http.StatusText(http.StatusOK)))

}
