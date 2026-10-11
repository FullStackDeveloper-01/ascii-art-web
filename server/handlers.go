package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"ascii-art-web/banner"
	"ascii-art-web/render"
)

const maxLength = 100

const maxBodyBytes = 1 << 20

var validBanners = map[string]bool{
	"standard":   true,
	"shadow":     true,
	"thinkertoy": true,
}

var bannerOrder = []string{"standard", "shadow", "thinkertoy"}

func recoverFromPanic(w http.ResponseWriter) {
	if rec := recover(); rec != nil {
		log.Println("panic recovered:", rec)
		render500(w)
	}
}

var static = http.StripPrefix("/static/", http.FileServer(http.Dir("static")))

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}
	static.ServeHTTP(w, r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverFromPanic(w)

	if r.URL.Path != "/" {
		render404(w, r)
		return
	}
	if r.Method != http.MethodGet {
		renderPage(w, http.StatusMethodNotAllowed, PageData{
			Banner: "standard",
			Error:  "405 Method Not Allowed: use GET",
		})
		return
	}

	renderPage(w, http.StatusOK, PageData{Banner: "standard"})
}

func ASCIIArtHandler(w http.ResponseWriter, r *http.Request) {
	defer recoverFromPanic(w)

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		renderPage(w, http.StatusMethodNotAllowed, PageData{
			Banner: "standard",
			Error:  "405 Method Not Allowed: use POST",
		})
		return
	}

	// ParseForm consumes the request body, so this endpoint intentionally
	// accepts POST only. GET requests are handled above with a 405 response.

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "400 Bad Request: invalid form data", http.StatusBadRequest)
		return
	}

	rawText := r.PostFormValue("text")
	bannerName := r.PostFormValue("banner")
	if bannerName == "" {
		bannerName = "standard"
	}
	action := r.PostFormValue("action")

	if !validBanners[bannerName] && bannerName != "all" {
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

	if count := utf8.RuneCountInString(text); count > maxLength {
		renderPage(w, http.StatusBadRequest, PageData{
			Text:   rawText,
			Banner: bannerName,
			Error:  fmt.Sprintf("Text is too long: %d characters (the limit is %d)", count, maxLength),
		})
		return
	}

	names := []string{bannerName}
	if bannerName == "all" {
		names = bannerOrder
	}

	var results []BannerResult
	for _, name := range names {
		chars, err := banner.Load(name)
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
		results = append(results, BannerResult{Name: name, Art: out})
	}

	if action == "download" {
		serveDownload(w, results)
		return
	}

	data := PageData{Text: rawText, Banner: bannerName}
	if bannerName == "all" {
		data.Results = results
	} else {
		data.Result = results[0].Art
	}
	renderPage(w, http.StatusOK, data)
}

func methodNotAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	http.Error(w, "405 Method Not Allowed: use "+allowed, http.StatusMethodNotAllowed)
}
