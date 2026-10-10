package main

import (
	"log"
	"net/http"
	"time"

	"ascii-art-web/server"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HomeHandler)
	mux.HandleFunc("/ascii-art", server.ASCIIArtHandler)
	mux.HandleFunc("/static/", server.StaticHandler)

	srv := &http.Server{
		Addr:         ":8090",
		Handler:      server.Logger(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Println("Server is running on http://localhost:8090")
	log.Fatal(srv.ListenAndServe())
}
