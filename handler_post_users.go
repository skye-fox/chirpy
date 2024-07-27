package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerPostUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.db.CreateUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
	}

	respondWithJSON(w, http.StatusCreated, user)

}
