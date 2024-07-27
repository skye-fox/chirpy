package main

import (
	"log"
	"net/http"

	"github.com/skye-fox/chirpy/internal/database"
)

type apiConfig struct {
	db             database.DB
	fileserverHits int
}

func main() {
	const (
		filepathRoot = "."
		port         = "8080"
		dbPath       = "database.json"
	)

	chirpDB, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	apiCFG := apiConfig{
		db:             *chirpDB,
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	fsHandler := apiCFG.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCFG.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCFG.handlerPostChirps)
	mux.HandleFunc("GET /api/chirps", apiCFG.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpid}", apiCFG.handlerGetChirpById)

	mux.HandleFunc("GET /admin/metrics", apiCFG.handlerMetrics)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
