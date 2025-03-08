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

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepath)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerHits)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	server := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}
	log.Printf("serving files from '%v', running on port: %v", filepath, port)
	log.Fatal(server.ListenAndServe())

}






