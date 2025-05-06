package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (config *apiConfig) apiUsersUpdateHandler(resp http.ResponseWriter, req *http.Request) {
	authorization, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithErrorJson(resp, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userId, err := auth.ValidateJWT(authorization, config.wjtSecret)
	if err != nil {
		respondWithErrorJson(resp, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	decoder := json.NewDecoder(req.Body)
	userInput := usersBody{}
	err = decoder.Decode(&userInput)
	if err != nil {
		respondWithErrorJson(resp, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashed, err := auth.HashPassword(userInput.Password)
	if err != nil {
		respondWithErrorJson(resp, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := config.db.UpdateUserCredentials(req.Context(), database.UpdateUserCredentialsParams{
		ID:             userId,
		Email:          userInput.Email,
		HashedPassword: hashed},
	)
	if err != nil {

		respondWithErrorJson(resp, http.StatusInternalServerError, "Error updating user credentials", err)
		return
	}

	respondWithJSON(resp, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
