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
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Println("Server is running on http://localhost" + port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
