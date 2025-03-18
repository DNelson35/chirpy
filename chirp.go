package main

import (
	"encoding/json"
	"net/http"
	"sort"
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

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		sendError(w, 500, "Couldn't retrieve chirps")
		return
	}

	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			sendError(w, 400, "Invalid author ID")
			return
		}
	}

	chirps := []respVal{}
	for _, dbChirp := range dbChirps {
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}

		chirps = append(chirps, respVal{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	sortOrder := r.URL.Query().Get("sort")

	switch sortOrder {
	case "asc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	case "desc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	sendOK(w, 200, &chirps)
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