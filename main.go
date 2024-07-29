package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/skye-fox/chirpy/internal/database"
)

type apiConfig struct {
	db             database.DB
	fileserverHits int
	jwtSecret      string
}

func main() {
	const (
		filepathRoot = "."
		port         = "8080"
		dbPath       = "database.json"
	)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("couldn't load environment variables")
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	chirpDB, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := chirpDB.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	apiCFG := apiConfig{
		db:             *chirpDB,
		fileserverHits: 0,
		jwtSecret:      jwtSecret,
	}

	mux := http.NewServeMux()
	fsHandler := apiCFG.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCFG.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCFG.handlerPostChirps)
	mux.HandleFunc("GET /api/chirps", apiCFG.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpid}", apiCFG.handlerGetChirpById)
	mux.HandleFunc("POST /api/login", apiCFG.handlerLogin)
	mux.HandleFunc("POST /api/users", apiCFG.handlerPostUsers)
	mux.HandleFunc("PUT /api/users", apiCFG.handlerUpdateUsers)

	mux.HandleFunc("GET /admin/metrics", apiCFG.handlerMetrics)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
