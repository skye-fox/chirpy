package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func (cfg *apiConfig) handlerPostChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if validateChirps(w, params.Body) {
		cleanedBody := cleanBody(params.Body)
		chirp, err := cfg.db.CreateChirp(cleanedBody)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		}

		respondWithJSON(w, http.StatusCreated, chirp)
	}
}

func validateChirps(w http.ResponseWriter, body string) bool {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return false
	}
	return true
}

func cleanBody(body string) string {
	const censored = "****"
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	for i, word := range words {
		if slices.Contains(bannedWords, strings.ToLower(word)) {
			words[i] = censored
		}
	}
	return strings.Join(words, " ")
}
