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

	apiCFG := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	fsHandler := apiCFG.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiCFG.handlerMetrics)
	mux.HandleFunc("GET /reset", apiCFG.handlerReset)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits int
}
