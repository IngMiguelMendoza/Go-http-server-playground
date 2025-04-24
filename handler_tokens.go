package main

import (
	"chirpy/internal/auth"
	"log"
	"net/http"
	"time"
)

func (config *apiConfig) apiTokenRefreshHandler(resp http.ResponseWriter, req *http.Request) {
	type NewToken struct {
		Token string `json:"token"`
	}

	inputToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithErrorJson(resp, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	refreshToken, err := config.db.GetRefreshToken(req.Context(), inputToken)
	if err != nil || time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
		respondWithErrorJson(resp, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	token, err := auth.MakeJWT(refreshToken.UserID, config.wjtSecret, time.Hour)
	if err != nil {
		respondWithErrorJson(resp, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}
	respondWithJSON(resp, http.StatusOK, NewToken{
		Token: token,
	})
}

func (config *apiConfig) apiTokenRevokeHandler(resp http.ResponseWriter, req *http.Request) {
	inputToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithErrorJson(resp, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	refreshToken, err := config.db.GetRefreshToken(req.Context(), inputToken)
	if err != nil {
		respondWithErrorJson(resp, http.StatusBadRequest, "Couldn't find stored token", err)
		return
	}

	err = config.db.RevokeToken(req.Context(), refreshToken.Token)
	if err != nil {
		log.Printf("Cannot revoke refresh token, %s", err)
		respondWithErrorJson(resp, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	respondWithJSON(resp, http.StatusNoContent, nil)
}
