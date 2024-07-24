package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const (
		filepathRoot = "."
		port         = "8080"
	)

	apiCFG := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCFG.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCFG.handlerMetrics)
	mux.HandleFunc("/reset", apiCFG.handlerReset)

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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	hits := fmt.Sprintf("Hits: %d", cfg.fileserverHits)
	w.Header().Add("Content-Type", "text/plain charset=utf-8")
	w.Write([]byte(hits))
}
