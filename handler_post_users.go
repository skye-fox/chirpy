package main

import (
	"encoding/json"
	"net/http"

	"github.com/skye-fox/chirpy/internal/auth"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerPostUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	if !cfg.db.CheckDuplicateEmail(params.Email) {
		user, err := cfg.db.CreateUser(params.Email, hashedPassword)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
			return
		}

		respondWithJSON(w, http.StatusCreated, User{
			Id:    user.Id,
			Email: user.Email,
		})
		return
	}

	respondWithError(w, http.StatusOK, "Account already exists")
}
