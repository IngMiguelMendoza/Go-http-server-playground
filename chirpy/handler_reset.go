package main

import (
	"net/http"
)

func (config *apiConfig) resetMetricsHandler(resp http.ResponseWriter, req *http.Request) {
	if config.platform != "dev" {
		resp.WriteHeader(http.StatusForbidden)
		resp.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	config.fileserverHits.Store(0)
	config.db.Reset(req.Context())
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
