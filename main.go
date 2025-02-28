package main

import (
	"log"
	"net/http"
)



func main() {
	port := "8080"
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	log.Printf("running on port: %v", port)
	log.Fatal(server.ListenAndServe())

}