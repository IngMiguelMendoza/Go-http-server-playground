package main

import (
	"chirpy/internal/auth"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (config *apiConfig) apiDeleteChirpByIdHandler(resp http.ResponseWriter, req *http.Request) {
	authorization, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithErrorJson(resp, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := auth.ValidateJWT(authorization, config.wjtSecret)
	if err != nil {
		respondWithErrorJson(resp, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	reqId, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithErrorJson(resp, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}

	record, err := config.db.GetChirp(req.Context(), reqId)
	if err != nil {
		log.Printf("Error retriving database record, %s", err)
		respondWithErrorJson(resp, http.StatusNotFound, "Error retriving database record", err)
		return
	}

	if userId != record.UserID {
		respondWithErrorJson(resp, http.StatusForbidden, "You can't delete this chirp", err)
		return
	}

	err = config.db.DeleteChirp(req.Context(), record.ID)
	if err != nil {
		log.Printf("Error retriving database record, %s", err)
		respondWithErrorJson(resp, http.StatusNotFound, "Error retriving database record", err)
		return
	}

	respondWithJSON(resp, http.StatusNoContent, nil)
}
