package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/skye-fox/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		Id           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.db.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusOK, "Couldn't find user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	expiration := 3600
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = expiration
	} else if params.ExpiresInSeconds > 3600 {
		params.ExpiresInSeconds = expiration
	}

	token, err := auth.MakeJWT(user.Id, cfg.jwtSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't Create JWT")
		return
	}

	var refreshToken string
	if user.RefreshToken == "" {
		refreshToken, err = auth.GenerateRefreshToken()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't Create refresh token")
		}

		cfg.db.AddRefreshToken(refreshToken, user.Id)
	}

	respondWithJSON(w, http.StatusOK, response{
		Id:           user.Id,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
