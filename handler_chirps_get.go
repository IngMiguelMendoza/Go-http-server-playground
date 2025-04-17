package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (config *apiConfig) apiGetChirpHandler(resp http.ResponseWriter, req *http.Request) {
	records, err := config.db.ListChirps(req.Context())
	if err != nil {
		log.Printf("Error retriving database records, %s", err)
		respondWithErrorJson(resp, 500, "Error retriving database records, %s", err)
		return
	}

	chirps := []ChirpRecord{}
	for _, record := range records {
		chirps = append(chirps, ChirpRecord{
			ID:        record.ID,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
			Body:      record.Body,
			UserID:    record.UserID,
		})
	}

	respondWithJSON(resp, http.StatusOK, chirps)
}

func (config *apiConfig) apiGetChirpByIdHandler(resp http.ResponseWriter, req *http.Request) {
	reqId, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("Invalid chirp id, %s", err)
		respondWithErrorJson(resp, http.StatusBadRequest, "Invalid chirp id, %s", err)
		return
	}

	record, err := config.db.GetChirp(req.Context(), reqId)
	if err != nil {
		log.Printf("Error retriving database record, %s", err)
		respondWithErrorJson(resp, http.StatusNotFound, "Error retriving database record, %s", err)
		return
	}

	respondWithJSON(resp, http.StatusOK, ChirpRecord{
		ID:        record.ID,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		Body:      record.Body,
		UserID:    record.UserID,
	})
}
