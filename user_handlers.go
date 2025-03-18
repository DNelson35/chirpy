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
	IsChirpyRed bool `json:"is_chirpy_red"`
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
		IsChirpyRed: user.IsChirpyRed,
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
		IsChirpyRed: user.IsChirpyRed,
	}

	sendOK(w, 200, &resp_user)
}

func(cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request){
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendError(w, 401, "unauthorized failed to get token")
		return
	}
	userID, err := auth.ValidateJWT(tokenStr, cfg.secretKey)
	if err != nil {
		sendError(w, 401, "unauthorized token invalid")
	}
	decoder := json.NewDecoder(r.Body)
	var req params

	if err = decoder.Decode(&req); err != nil {
		sendError(w, 400, "bad request")
		return
	}

	hashPas, err := auth.HashPassword(req.Password)
	if err != nil {
		sendError(w, 500, "failed to hash password")
		return
	}
	user, err := cfg.db.UpdateUserData(r.Context(), database.UpdateUserDataParams{
		Email: req.Email,
		HashedPassword: hashPas,
		ID: userID,
	})
	if err != nil {
		sendError(w, 500, "failed to update user")
	}

	resp := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	sendOK(w, 200, &resp)
}

func(cfg *apiConfig) handlePolkaWebHook(w http.ResponseWriter, r *http.Request){
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
	err := cfg.db.UpgradeUserChirpyRed(r.Context(), req.Data.UserID)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
}

