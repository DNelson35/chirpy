package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DNelson35/chirpy/internal/auth"
	"github.com/DNelson35/chirpy/internal/database"
	"github.com/google/uuid"
)

type params struct {
	Email string `json:"email"`
	Password string `json:"password"`
	ExpiresIn int `json:"expires_in_seconds"`
}

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func(cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var resp params
	if err := decoder.Decode(&resp); err != nil {
		sendError(w, 500, "Something went Wrong")
		return
	}

	pass, err := auth.HashPassword(resp.Password)
	if err != nil{
		sendError(w, 500, "Failed to hash password")
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: resp.Email,
		HashedPassword: pass,
	})

	if err != nil {
		sendError(w, 400, "Failed to create user")
		return
	}

	resp_user := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	sendOK(w, 201, &resp_user)
}

func(cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var resp params
	if err := decoder.Decode(&resp); err != nil {
		sendError(w, 500, "Something went Wrong")
		return
	}
	user, err := cfg.db.GetUser(r.Context(), resp.Email)
	if err != nil {
		sendError(w, 404, "User not found")
		return 
	}
	
	if err = auth.CheckPassword(resp.Password, user.HashedPassword); err != nil {
		sendError(w, 401, "Password incorrect not authorized")
		return
	}

	if resp.ExpiresIn == 0 || resp.ExpiresIn > 3600 {
		resp.ExpiresIn = 3600
	}

	tokenString, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Duration(resp.ExpiresIn) * time.Second)
	if err != nil {
		sendError(w, 401, "Unauthorized")
		return
	}
	refreshTokenStr := auth.MakeRefreshToken()
	refToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshTokenStr,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		sendError(w, 400, "failed to create refresh token")
		return 
	}
	
	resp_user := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: tokenString,
		RefreshToken: refToken.Token,
	}

	sendOK(w, 200, &resp_user)
}

