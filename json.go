package main

import "encoding/json"
import "net/http"

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorRes struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorRes{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
	return
}
