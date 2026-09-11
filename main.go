package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"ascii-art-web/banner"
	"ascii-art-web/render"
)

// maxLength mirrors the limit used by the original CLI project.
const maxLength = 10000

const templatePath = "templates/index.html"

// Change this if 8090 is ever busy on your machine — just this one line.
const port = ":8090"

// validBanners is the whitelist of banner names the form is allowed to send.
var validBanners = map[string]bool{
	"standard":   true,
	"shadow":     true,
	"thinkertoy": true,
}

// PageData is everything the HTML template needs. Empty fields are fine:
// the template shows/hides sections based on them (see {{if .Error}} etc).
type PageData struct {
	Text   string
	Banner string
	Result string
	Error  string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/ascii-art", asciiArtHandler)

	// Serve style.css (and anything else placed in static/) at /static/.
	// http.FileServer already returns 404 by itself for a missing file.

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Println("Server is running at http://localhost" + port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

// recoverFromPanic turns any unexpected panic into a 500 response instead of
// crashing the whole server. Call it with `defer recoverFromPanic(w)` as the
// very first line of a handler.

func recoverFromPanic(w http.ResponseWriter) {
	if rec := recover(); rec != nil {
		log.Println("panic recovered:", rec)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

// homeHandler serves GET /. Any path other than "/" is a 404 because
// http.ServeMux routes every unmatched path to the "/" handler by default.

func homeHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverFromPanic(w)

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "400 Bad Request: use GET", http.StatusBadRequest)
		return
	}

	renderPage(w, http.StatusOK, PageData{Banner: "standard"})
}

// asciiArtHandler serves POST /ascii-art. It validates the input, calls the
// exact same banner.Load / render.Render functions the CLI project uses,
// and turns any error into the appropriate HTTP status code.

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverFromPanic(w)

	if r.Method != http.MethodPost {
		http.Error(w, "400 Bad Request: use POST", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "400 Bad Request: invalid form data", http.StatusBadRequest)
		return
	}

	rawText := r.FormValue("text")
	bannerName := r.FormValue("banner")
	if bannerName == "" {
		bannerName = "standard"
	}

	if !validBanners[bannerName] {
		renderPage(w, http.StatusBadRequest, PageData{
			Text:   rawText,
			Banner: bannerName,
			Error:  "Unknown banner: " + bannerName,
		})
		return
	}

	// Same trick the CLI used: a literal backslash-n typed by the user
	// becomes a real newline, so "Hello\nThere" still works from a
	// single-line <input>.

	text := strings.ReplaceAll(rawText, `\n`, "\n")

	if len(text) > maxLength {
		renderPage(w, http.StatusBadRequest, PageData{
			Text:   rawText,
			Banner: bannerName,
			Error:  "Text is too long",
		})
		return
	}

	chars, err := banner.Load(bannerName)
	if err != nil {
		// bannerName is already whitelisted above, so this should not
		// normally happen — handled anyway to stay safe.
		renderPage(w, http.StatusBadRequest, PageData{
			Text:   rawText,
			Banner: bannerName,
			Error:  err.Error(),
		})
		return
	}

	out, err := render.Render(text, chars)
	if err != nil {
		// e.g. "character not supported" for non-ASCII / control chars.
		renderPage(w, http.StatusBadRequest, PageData{
			Text:   rawText,
			Banner: bannerName,
			Error:  err.Error(),
		})
		return
	}

	renderPage(w, http.StatusOK, PageData{
		Text:   rawText,
		Banner: bannerName,
		Result: out,
	})
}

// renderPage parses the template fresh on every call. If the template file
// is missing, that's a 404 (a required resource wasn't found), not a 500.
func renderPage(w http.ResponseWriter, status int, data PageData) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "404 Not Found: template missing", http.StatusNotFound)
		return
	}

	w.WriteHeader(status)
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("template execute error:", err)
	}
}
