package main

import (
	"fmt"
	"net/http"
)

func (config *apiConfig) metricsIncrement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		config.fileserverHits.Add(1)
		next.ServeHTTP(resp, req)
	})
}

func (config *apiConfig) metricsHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(fmt.Sprintf(`<html>
<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>
</html>`, config.fileserverHits.Load())))
}
