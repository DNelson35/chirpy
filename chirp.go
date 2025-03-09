package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/DNelson35/chirpy/internal/database"
	"github.com/google/uuid"
)



type respVal struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type reqVal struct {
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func(cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)

	var req reqVal

	if err := decoder.Decode(&req); err != nil {
		sendError(w, 500, "Something went wrong")
		return
	}

	if len(req.Body) > 140 {
		sendError(w, 400, "Chirp is too long")
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanInput(&req),
		UserID: req.UserID,
	})
	if err != nil {
		sendError(w, 400, "failed to create chirp")
	}

	resp := respVal{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	sendOK(w, 201, &resp)
	return
}



func cleanInput(req *reqVal) string {
	profainWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	wordList := strings.Split(req.Body, " ")

	for i, word := range wordList{
		if _, ok := profainWords[strings.ToLower(word)]; ok {
			wordList[i] = "****"
		}
	}
	return strings.Join(wordList, " ")
}