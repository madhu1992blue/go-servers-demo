package main

import (
	"time"
	"github.com/google/uuid"
	"net/http"
	"encoding/json"
	"github.com/madhu1992blue/go-servers-demo/internal/database"
	"github.com/madhu1992blue/go-servers-demo/internal/auth"
)
type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(),params.Email)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword) ; err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	respondWithJSON(w, 200, User{
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ID: user.ID,
	})
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var params parameters

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, "Bad request")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	newUser := User {
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	}
	respondWithJSON(w, 201, newUser)
}
