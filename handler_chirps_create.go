package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ChirpInput struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

type ChirpRecord struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (config *apiConfig) apiPostChirpHandler(resp http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirpInput := ChirpInput{}
	err := decoder.Decode(&chirpInput)
	if err != nil {
		respondWithErrorJson(resp, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	authorization, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error, %s, for user %v request, ", err, chirpInput.UserId.String())
		respondWithErrorJson(resp, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userId, err := auth.ValidateJWT(authorization, config.wjtSecret)
	if err != nil {
		log.Printf("Unauthorized request, %s", err)
		respondWithErrorJson(resp, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// Original content is loss, norequirement to sve the original content yet
	err = validateChirp(&chirpInput)
	if err != nil {
		log.Printf("Error chirp body content too big, length %d", len(chirpInput.Body))
		respondWithErrorJson(resp, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := config.db.CreateChirps(req.Context(), database.CreateChirpsParams{
		Body:   chirpInput.Body,
		UserID: userId},
	)
	if err != nil {
		log.Printf("Error decoding request body, %s", err)
		respondWithErrorJson(resp, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(resp, http.StatusCreated, ChirpRecord{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
