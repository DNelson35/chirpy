package main

import (
	"net/http"
	"fmt"
)

func(cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func(cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := cfg.fileServerHits.Load()
	w.Write([]byte(fmt.Sprintf(`
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
	`, hits)))
}