package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)


type apiConfig struct {
	fileServerHits atomic.Int32
}

func(cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := "8080"
	filepath := "."

	var cfg apiConfig

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepath)))))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", cfg.handlerHits)
	mux.HandleFunc("/reset", cfg.handlerReset)

	server := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}
	log.Printf("serving files from '%v', running on port: %v", filepath, port)
	log.Fatal(server.ListenAndServe())

}

func handlerReadiness(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func(cfg *apiConfig) handlerHits (w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileServerHits.Load()
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

func(cfg *apiConfig) handlerReset (w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileServerHits.Store(0)
	w.Write([]byte(fmt.Sprintf("Hits Reset: %d", cfg.fileServerHits.Load())))
}