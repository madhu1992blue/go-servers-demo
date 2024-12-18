package main

import "net/http"
import "fmt"

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
func main() {
	mux := &http.ServeMux{}
	mux.Handle("/app/", http.StripPrefix("/app/",http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", healthz)
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("error occured: %v", err)
	}
}
