package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/DNelson35/chirpy/internal/auth"
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
}

func(cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request){
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendError(w, 401, "Unauthorized")
	}

	userId, err := auth.ValidateJWT(tokenString, cfg.secretKey)
	if err != nil {
		sendError(w, 401, "Unauthorized")
	}
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
		UserID: userId,
	})
	if err != nil {
		sendError(w, 400, "failed to create chirp")
		return
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

func(cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request){
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		sendError(w, 400, "could not get chirps")
		return
	}
	var chirpsList []respVal
	for _, chirp := range chirps {
		newChirp := respVal{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		}
		chirpsList = append(chirpsList, newChirp)
	}

	sendOK(w, 200, &chirpsList)

}

func(cfg *apiConfig) handlerGetChirpsById(w http.ResponseWriter, r *http.Request){
	path , err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		sendError(w, 400, "Invalid ID format")
		return
	}
	chirp, err := cfg.db.GetChirpsById(r.Context(), path)

	if err != nil {
		sendError(w, 404, "Chirp not found")
		return
	}

	resp := respVal{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	sendOK(w, 200, &resp)
}

func(cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request){
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendError(w, 401, "unauthorized failed to get token")
		return
	}
	userID, err := auth.ValidateJWT(tokenStr, cfg.secretKey)
	if err != nil {
		sendError(w, 401, "unauthorized token invalid")
		return
	}
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		sendError(w,400,"invalid id format")
		return
	}
	chirp, err := cfg.db.GetChirpsById(r.Context(), chirpID)
	if err != nil {
		sendError(w, 404, "chirp not found")
	}
	if chirp.UserID != userID {
		sendError(w, 403, "Delete Not allowed")
		return
	}
	err = cfg.db.DeleteChipById(r.Context(), chirp.ID)
	if err != nil {
		sendError(w, 500, "failed to delete chirp")
	}
	w.WriteHeader(204)
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