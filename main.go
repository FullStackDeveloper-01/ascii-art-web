package main

import (
	"log"
	"net/http"
	"time"

	"ascii-art-web/server"
)

const port = ":8090"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HomeHandler)
	mux.HandleFunc("/ascii-art", server.ASCIIArtHandler)

	fileServer := http.FileServer(http.Dir("static"))
	staticHandler := http.StripPrefix("/static/", fileServer)
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "400 Bad Request: use GET", http.StatusBadRequest)
			return
		}
		staticHandler.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr:              port,
		Handler:           server.Logger(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("Server is running on http://localhost" + port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
