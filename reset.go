package main

import (
	"fmt"
	"net/http"
	"os"
)

func(cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	dev := os.Getenv("PLATFORM")
	if dev != "dev" {
		sendError(w, 403, "Forbidden")
		return
	}
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		sendError(w, 500, "failed to delete users")
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileServerHits.Store(0)
	w.Write([]byte(fmt.Sprintf("Hits Reset: %d", cfg.fileServerHits.Load())))
}