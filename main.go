package main

import (
	"log"
	"net/http"
	"sync/atomic"
)


type apiConfig struct {
	fileServerHits atomic.Int32
}

func main() {
	port := "8080"
	filepath := "."

	var cfg apiConfig

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






