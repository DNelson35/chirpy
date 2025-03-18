package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/DNelson35/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


type apiConfig struct {
	fileServerHits atomic.Int32
	db *database.Queries
	secretKey string
}

func main() {
	port := "8080"
	filepath := "."
	
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		fmt.Println("failed to connect")
	}
	dbQueries := database.New(db)

	var cfg apiConfig
	cfg.db = dbQueries
	cfg.secretKey = os.Getenv("SECRET")

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepath)))))
	mux.HandleFunc("GET /admin/metrics", cfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirpsById)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handleDeleteChirp)
	mux.HandleFunc("POST /api/users", cfg.handleCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.handleUpdateUser)
	mux.HandleFunc("POST /api/login", cfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.handleRevokeToken)


	server := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}
	log.Printf("serving files from '%v', running on port: %v", filepath, port)
	log.Fatal(server.ListenAndServe())

}






