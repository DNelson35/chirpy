package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func(cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email string `json:"email"`
	}
	type User struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	var resp params
	if err := decoder.Decode(&resp); err != nil {
		sendError(w, 500, "Something went Wrong")
	}

	user, err := cfg.db.CreateUser(r.Context(), resp.Email)
	if err != nil {
		sendError(w, 400, "Failed to create user")
	}

	resp_user := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	sendOK(w, 201, &resp_user)
}

