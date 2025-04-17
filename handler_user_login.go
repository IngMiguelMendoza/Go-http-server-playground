package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	MAX_EXPIRATION_TIME int = 60 * 60
)

func (config *apiConfig) apiUserLoginHandler(resp http.ResponseWriter, req *http.Request) {
	type inputCredentials struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(req.Body)
	credentials := inputCredentials{}
	err := decoder.Decode(&credentials)
	if err != nil {
		respondWithErrorJson(resp, http.StatusInternalServerError, "Error decoding login request body", err)
		return
	}

	user, err := config.db.GeUserByMail(req.Context(), credentials.Email)
	if err != nil {
		respondWithErrorJson(resp, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(credentials.Password, user.HashedPassword)
	if err != nil {
		respondWithErrorJson(resp, http.StatusUnauthorized, "Password mismatch, %s", err)
		return
	}

	expirationTime := time.Hour
	if credentials.ExpiresInSeconds > 0 && credentials.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(credentials.ExpiresInSeconds) * time.Second
	}

	token, err := auth.MakeJWT(user.ID, config.wjtSecret, expirationTime)
	if err != nil {
		log.Printf("Cannot create the JWT token, %s", err)
	}
	respondWithJSON(resp, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
