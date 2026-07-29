package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bePramudya/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type returnVals struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{
		ExpiresInSeconds: 60 * 60,
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't compare passwrord hash", err)
		return
	}
	if match == false {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	expirationTime := time.Duration(params.ExpiresInSeconds) * time.Second
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		User: User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: accessToken,
	})

}
