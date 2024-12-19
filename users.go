package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/madhu1992blue/go-servers-demo/internal/auth"
	"github.com/madhu1992blue/go-servers-demo/internal/database"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
	}

	type UserLoggedIn struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	params := &parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(params); err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	tokenString, err := auth.MakeJWT(user.ID, cfg.jwtSecret, 60 * 60 * time.Second)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}
	refreshToken, err := auth.MakeRefreshToken()
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(60 * 60 * 24 * time.Second),
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}
	respondWithJSON(w, 200, UserLoggedIn{
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ID:        user.ID,
		Token:     tokenString,
		RefreshToken: refreshToken,
	})
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
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
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	newUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJSON(w, 201, newUser)
}
