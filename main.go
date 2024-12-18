package main

import "net/http"
import "fmt"
func main() {
	mux := &http.ServeMux{}
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("error occured: %v", err)
	}
}
