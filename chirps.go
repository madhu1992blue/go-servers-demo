package main
import (
	"net/http"
	"github.com/madhu1992blue/go-servers-demo/internal/database"
	"github.com/madhu1992blue/go-servers-demo/internal/auth"
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

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Println("Error: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	chirpData, err := cfg.db.GetChirp(r.Context(), chirpID);
	if err != nil {
		log.Println("Error: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	chirp := Chirp {
		Body: chirpData.Body,
		CreatedAt: chirpData.CreatedAt,
		UpdatedAt: chirpData.UpdatedAt,
		ID: chirpData.ID,
		UserID: chirpData.UserID,
	}
	respondWithJSON(w, 200, chirp)
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
	tokenString, err := auth.GetBearerToken(&r.Header)
	if err != nil {
		log.Printf("Error : %v", err)
		respondWithError(w, 401, "Unauthorized")
		return
	}
	userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret) 
	if err!=nil {
		log.Printf("Error : %v", err)
		respondWithError(w, 401, "Unauthorized")
	}
	type parameters struct {
		Body string `json:"body"`
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
		UserID: userID,
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
