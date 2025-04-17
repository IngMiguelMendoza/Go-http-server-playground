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

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (config *apiConfig) usersHandler(resp http.ResponseWriter, req *http.Request) {
	type usersBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	userInput := usersBody{}
	err := decoder.Decode(&userInput)
	if err != nil {
		log.Printf("Error decoding request body, %s", err)
		respondWithErrorJson(resp, 400, "Invalid user mail record", err)
		return
	}

	hashed, err := auth.HashPassword(userInput.Password)

	user, err := config.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          userInput.Email,
		HashedPassword: hashed,
	})

	if err != nil {
		log.Printf("Error decoding request body, %s", err)
		respondWithErrorJson(resp, 500, "Could not create user record on the databse, %s", err)
		return
	}

	respondWithJSON(resp, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
