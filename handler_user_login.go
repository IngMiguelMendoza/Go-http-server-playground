package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"
)

const (
	MAX_EXPIRATION_TIME           int = 60 * 60
	REFRESH_TOKEN_EXPIRATION_TIME int = 60 * 24
)

func (config *apiConfig) apiUserLoginHandler(resp http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(req.Body)
	credentials := parameters{}
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

	accessToken, err := auth.MakeJWT(user.ID, config.wjtSecret, time.Hour)
	if err != nil {
		respondWithErrorJson(resp, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	stringToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithErrorJson(resp, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	refreshToken, err := config.db.CreateToken(req.Context(), database.CreateTokenParams{
		Token:     stringToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(REFRESH_TOKEN_EXPIRATION_TIME) * time.Hour),
	})
	if err != nil {
		respondWithErrorJson(resp, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	respondWithJSON(resp, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken.Token,
	})
}
