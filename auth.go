package main
import (
	"time"
	"net/http"
	"github.com/madhu1992blue/go-servers-demo/internal/auth"
)
func (cfg *apiConfig) RefreshAccessToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(&r.Header)
	if err != nil {
		respondWithError(w, 400, "Bad request")
		return
	}
	tokenRecord, err := cfg.db.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	if tokenRecord.ExpiresAt.Before(time.Now()) || tokenRecord.RevokedAt.Valid {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	signedJWT, err := auth.MakeJWT(tokenRecord.UserID, cfg.jwtSecret, 60 * 60 * time.Second)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	result := struct { Token string `json:"token"` }{
		Token: signedJWT,
	}
	respondWithJSON(w, 200, result)
}


func (cfg *apiConfig) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err:= auth.GetBearerToken(&r.Header)
	if err != nil {
		respondWithError(w, 401, "Unauthorized")
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(),refreshToken)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	respondWithJSON(w,204,nil)
}
