package main
import (
	"github.com/madhu1992blue/go-servers-demo/internal/database"
	"net/http"
	"github.com/google/uuid"
	"encoding/json"
	"strings"
	"log"
	"time"
)

type Chirp struct {
	Body string `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirpsData, err := cfg.db.GetChirps(r.Context());
	if err != nil {
		log.Println("Error: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	chirps := make([]Chirp, len(chirpsData))
	for i, c := range chirpsData {
		chirps[i].Body = c.Body
		chirps[i].CreatedAt = c.CreatedAt
		chirps[i].UpdatedAt = c.UpdatedAt
		chirps[i].ID = c.ID
		chirps[i].UserID = c.UserID
	}
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type errorRes struct {
		Error string `json:"error"`
	}
	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	chirpParts := strings.Fields(params.Body)
	for i, part := range chirpParts {
		for _, pw := range profaneWords {
			if strings.ToLower(part) == strings.ToLower(pw) {
				chirpParts[i] = "****"
				continue
			}
		}
	}
	cleanedBody := strings.Join(chirpParts, " ")
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedBody,
		UserID: params.UserID,
	})
	if err != nil {
		log.Println("Error: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	respondWithJSON(w, 201, Chirp{
		ID: chirp.ID,
		Body: chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID: chirp.UserID,
	})
}
