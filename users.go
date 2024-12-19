package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/madhu1992blue/go-servers-demo/internal/auth"
	"github.com/madhu1992blue/go-servers-demo/internal/database"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsChirpyRed bool	`json:"is_chirpy_red"`
}
type UserLoggedIn struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool	`json:"is_chirpy_red"`
}
func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event    string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}
	params := &parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(params); err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}
	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}
	user, err := cfg.db.GetUserByID(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "record not found")
		return
	}
	err = cfg.db.UpgradeUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, "Somethign went wrong")
		return
	}
	respondWithJSON(w, 204, nil)

}
func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := &parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(params); err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}
	accessTokenString, err := auth.GetBearerToken(&r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(accessTokenString, cfg.jwtSecret)
	if err != nil {
		log.Printf("JWT Validation ERROR: %v -- %s", err, userID)
		respondWithError(w, 401, "Unauthorized")
		return
	}
	_, err = cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("ERROR: %v", err)
		respondWithError(w, 401, "Unauthorized")
		return
	}
	err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("ERROR: %v", err)
		respondWithError(w, 501, "Something went wrong")
		return
	}
	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, 501, "Something went wrong")
		return
	}
	result := struct {
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		ID        uuid.UUID `json:"id"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}{
		Email:     params.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ID:        user.ID,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, 200, result)

}
func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	tokenString, err := auth.MakeJWT(user.ID, cfg.jwtSecret, 60*60*time.Second)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}
	refreshToken, err := auth.MakeRefreshToken()
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 60 * 24 * time.Second * 60),
	})
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}
	respondWithJSON(w, 200, UserLoggedIn{
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		ID:           user.ID,
		Token:        tokenString,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
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
		IsChirpyRed: dbUser.IsChirpyRed,
	}
	respondWithJSON(w, 201, newUser)
}
