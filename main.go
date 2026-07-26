package main

import (
	"log"
	"net/http"
)

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: &http.ServeMux{},
	}
	log.Fatal(server.ListenAndServe())
}
