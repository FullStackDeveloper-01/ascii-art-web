package main

import (
	"log"
	"net/http"

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

	log.Println("Server is running on http://localhost" + port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
