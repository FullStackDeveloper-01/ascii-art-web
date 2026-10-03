package main

import (
	"log"
	"net/http"

	"ascii-art-web/server"
)

// Change this if 8090 is ever busy on your machine — just this one line.
const port = ":8090"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.HomeHandler)
	mux.HandleFunc("/ascii-art", server.ASCIIArtHandler)

	// Serve style.css (and anything else placed in static/) at /static/.
	// http.FileServer already returns 404 by itself for a missing file.

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Println("Server is running on http://localhost" + port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
