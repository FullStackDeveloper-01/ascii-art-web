package server

import (
	"log"
	"net/http"
	"strings"

	"ascii-art-web/banner"
	"ascii-art-web/render"
)

const maxLength = 10000

var validBanners = map[string]bool{
	"standard":   true,
	"shadow":     true,
	"thinkertoy": true,
}

func recoverFromPanic(w http.ResponseWriter) {
	if rec := recover(); rec != nil {
		log.Println("panic recovered:", rec)
		render500(w)
	}
}

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

	rawText := r.PostFormValue("text")
	bannerName := r.PostFormValue("banner")
	if bannerName == "" {
		bannerName = "standard"
	}

	if !validBanners[bannerName] {
		renderPage(w, http.StatusNotFound, PageData{
			Text:   rawText,
			Banner: bannerName,
			Error:  "Unknown banner: " + bannerName,
		})
		return
	}

	if rawText == "" {
		renderPage(w, http.StatusBadRequest, PageData{
			Banner: bannerName,
			Error:  "Text is required",
		})
		return
	}

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
		log.Println("banner load error:", err)
		render500(w)
		return
	}

	out, err := render.Render(text, chars)
	if err != nil {
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
