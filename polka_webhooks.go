package main

import (
	"encoding/json"
	"net/http"

	"github.com/DNelson35/chirpy/internal/auth"
	"github.com/google/uuid"
)

func(cfg *apiConfig) handlePolkaWebHook(w http.ResponseWriter, r *http.Request){
	key, err := auth.GetApiKey(r.Header)
	if err != nil || key != cfg.polkaKey {
		w.WriteHeader(401)
		return
	}
	
	type hookResp struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	var req hookResp
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		sendError(w, 400, "bad request")
	}
	if req.Event != "user.upgraded"{
		w.WriteHeader(204)
		return
	}
	err = cfg.db.UpgradeUserChirpyRed(r.Context(), req.Data.UserID)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
}

