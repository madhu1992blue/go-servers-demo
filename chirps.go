package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/madhu1992blue/go-servers-demo/internal/auth"
	"github.com/madhu1992blue/go-servers-demo/internal/database"
	"log"
	"net/http"
	"strings"
	"time"
)

type Chirp struct {
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Println("Error: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	chirpData, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		log.Println("Error: %v", err)
		respondWithError(w, 404, "Record not found")
		return
	}
	chirp := Chirp{
		Body:      chirpData.Body,
		CreatedAt: chirpData.CreatedAt,
		UpdatedAt: chirpData.UpdatedAt,
		ID:        chirpData.ID,
		UserID:    chirpData.UserID,
	}
	respondWithJSON(w, 200, chirp)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	authorIdQuery := r.URL.Query().Get("author_id")
	var chirpsData []database.Chirp
	var err error
	if authorIdQuery !="" {
		authorId, err := uuid.Parse(authorIdQuery)
		if err != nil {
			respondWithError(w, 500, "Something went wrong")
			return
		}
		chirpsData, err = cfg.db.GetChirpsByAuthor(r.Context(), authorId)
	} else {
		chirpsData, err = cfg.db.GetChirps(r.Context())
	}
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
	if err != nil {
		log.Printf("Error : %v", err)
		respondWithError(w, 401, "Unauthorized")
		return
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
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		log.Println("Error: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}
	respondWithJSON(w, 201, Chirp{
		ID:        chirp.ID,
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
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
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}
	chirpRecord, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}
	if chirpRecord.UserID != userID {
		respondWithError(w, 403, "Unauthorized")
		return
	}
	err = cfg.db.DeleteChirpByIDAndUser(r.Context(),database.DeleteChirpByIDAndUserParams{
		ID: chirpId,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, 403, "Unauthorized")
		return
	}
	respondWithJSON(w, 204, nil)

}
