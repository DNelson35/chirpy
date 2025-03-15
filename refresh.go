package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/DNelson35/chirpy/internal/auth"
	"github.com/DNelson35/chirpy/internal/database"
)

func(cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request){
	type resp struct {
		Token string `json:"token"`
	}

	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendError(w, 400, "bad request")
		return
	}

	refToken, err := cfg.db.GetRefToken(r.Context(), tokenStr)
	if err != nil {
		sendError(w, 401, "Not authorized")
		return
	}

	isExpired := refToken.ExpiresAt.Before(time.Now())
	if isExpired || refToken.RevokedAt.Valid {
		sendError(w, 401, "Token has expired")
		return
	}

	accessToken, err := auth.MakeJWT(refToken.UserID, cfg.secretKey, time.Hour)
	if err != nil {
		sendError(w, 500, "Failed to generate access token")
		return
	}

	token := resp{
		Token: accessToken,
	}
	sendOK(w, 200, &token)
}

func(cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request){
	tokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendError(w, 400, "bad request")
	}
	err = cfg.db.UpdateRefTokenRevocation(r.Context(), database.UpdateRefTokenRevocationParams{
		Token: tokenStr,
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		sendError(w, 500, "failed to revoke token")
	}
	w.WriteHeader(204)
	return
}