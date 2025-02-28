package main

import (
	"log"
	"net/http"
)



func main() {
	port := "8080"
	filepath := "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepath)))

	server := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	
	

	log.Printf("serving files from '%v', running on port: %v", filepath, port)
	log.Fatal(server.ListenAndServe())

}