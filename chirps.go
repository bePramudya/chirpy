package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bePramudya/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirps struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type returnValues struct {
		Chirps
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode json", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	chirpArgs := database.CreateChirpParams{
		Body:   cleaned,
		UserID: params.UserId,
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), chirpArgs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnValues{
		Chirps: Chirps{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		},
	})
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")

	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := badWords[lowerWord]; ok {
			words[i] = "****"
		}
	}

	cleaned := strings.Join(words, " ")
	return cleaned
}
