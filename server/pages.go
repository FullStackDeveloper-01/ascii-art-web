package server

import (
	"html/template"
	"log"
	"net/http"

	"ascii-art-web/banner"
	"ascii-art-web/render"
)

const templatePath = "templates/index.html"
const notFoundTemplatePath = "templates/404.html"

// PageData is everything the HTML template needs. Empty fields are fine:
// the template shows/hides sections based on them (see {{if .Error}} etc).
type PageData struct {
	Text   string
	Banner string
	Result string
	Error  string
	Sample string
}

// renderPage parses the template fresh on every call. If the template file
// is missing, that's a 404 (a required resource wasn't found), not a 500.
func renderPage(w http.ResponseWriter, status int, data PageData) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "404 Not Found: template missing", http.StatusNotFound)
		return
	}

	data.Sample = generateSample()

	w.WriteHeader(status)
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("template execute error:", err)
	}
}

// generateSample renders a small fixed phrase through the exact same
// banner.Load/render.Render pipeline as every real submission, purely as a
// decorative preview in the page header — it does not reimplement or
// duplicate any conversion logic. If it ever fails for any reason, it just
// renders as an empty string (the header simply omits it); it never affects
// the response's status code.
func generateSample() string {
	chars, err := banner.Load("standard")
	if err != nil {
		return ""
	}
	out, err := render.Render("01", chars)
	if err != nil {
		return ""
	}
	return out
}

// errorPageData is the (tiny) data the 404 page needs.
type errorPageData struct {
	Sample string
}

// render404 serves the styled 404 page instead of Go's plain default text.
// The status code is still exactly 404 — only the page's look changes. If
// the 404 template itself can't be parsed for any reason, it falls back to
// the plain http.NotFound so a broken template file never turns into a
// worse error than a missing page.
func render404(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(notFoundTemplatePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	data := errorPageData{Sample: generate404Sample()}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("404 template execute error:", err)
	}
}

// generate404Sample renders "404" through the exact same banner pipeline as
// generateSample above — same safety guarantee: any failure just means the
// decoration is skipped, never a different status code.
func generate404Sample() string {
	chars, err := banner.Load("standard")
	if err != nil {
		return ""
	}
	out, err := render.Render("404", chars)
	if err != nil {
		return ""
	}
	return out
}
