package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/skye-fox/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	user, err := cfg.db.UpdateUsers(id, params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusOK, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		Id:    user.Id,
		Email: user.Email,
	})
}
