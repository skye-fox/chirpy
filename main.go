package main

import (
	"log"
	"net/http"
)

func main() {
	const (
		filepathRoot = "."
		port         = "8080"
	)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	server := http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
