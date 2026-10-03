package server

import (
	"log"
	"net/http"
	"strings"

	"ascii-art-web/banner"
	"ascii-art-web/render"
)

// maxLength mirrors the limit used by the original CLI project.
const maxLength = 10000

// validBanners is the whitelist of banner names the form is allowed to send.
var validBanners = map[string]bool{
	"standard":   true,
	"shadow":     true,
	"thinkertoy": true,
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

// HomeHandler serves GET /. Any path other than "/" is a 404 because
// http.ServeMux routes every unmatched path to the "/" handler by default.

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverFromPanic(w)

	if r.URL.Path != "/" {
		render404(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "400 Bad Request: use GET", http.StatusBadRequest)
		return
	}

	renderPage(w, http.StatusOK, PageData{Banner: "standard"})
}

// ASCIIArtHandler serves POST /ascii-art. It validates the input, calls the
// exact same banner.Load / render.Render functions the CLI project uses,
// and turns any error into the appropriate HTTP status code.

func ASCIIArtHandler(w http.ResponseWriter, r *http.Request) {
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
		// The subject lists "banners" explicitly as a 404 case (same bucket
		// as a missing template), not 400 — a banner that doesn't exist is a
		// "not found" resource, not a malformed request.
		renderPage(w, http.StatusNotFound, PageData{
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
		// bannerName is already whitelisted above, so this only happens if
		// the embedded banner data itself is broken — that's the server's
		// fault, not the visitor's request, so it's a 500, not a 400/404.
		log.Println("banner load error:", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
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
