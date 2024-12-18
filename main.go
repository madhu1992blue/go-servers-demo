package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"sync/atomic"
	"strings"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

var profaneWords []string = []string{"kerfuffle", "sharbert", "fornax"}
func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type errorRes struct {
		Error string `json:"error"`
	}
	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
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
	respondWithJSON(w, 200, returnVal{CleanedBody: cleanedBody})
}
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		w.Header().Add("Cache-Control", "no-cache")
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	messageFormat := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	w.Write([]byte(fmt.Sprintf(messageFormat,cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
}
func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
func main() {
	mux := &http.ServeMux{}
	cfg := &apiConfig{}
	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileServerHandler))
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("GET /admin/metrics", cfg.metrics)
	mux.HandleFunc("POST /admin/reset", cfg.reset)
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("error occured: %v", err)
	}
}
