package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"ascii-art-web/banner"
	"ascii-art-web/render"
)

const templatePath = "templates/index.html"
const notFoundTemplatePath = "templates/404.html"
const internalErrorTemplatePath = "templates/500.html"

type PageData struct {
	Text    string
	Banner  string
	Result  string
	Results []BannerResult
	Error   string
	Sample  string
}

type BannerResult struct {
	Name string
	Art  string
}

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

type errorPageData struct {
	Sample string
}

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

func render500(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles(internalErrorTemplatePath)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	data := errorPageData{Sample: generate500Sample()}
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("500 template execute error:", err)
	}
}

func generate500Sample() string {
	chars, err := banner.Load("standard")
	if err != nil {
		return ""
	}
	out, err := render.Render("500", chars)
	if err != nil {
		return ""
	}
	return out
}

func serveDownload(w http.ResponseWriter, results []BannerResult) {
	var body strings.Builder
	for i, res := range results {
		if len(results) > 1 {
			if i > 0 {
				body.WriteString("\n")
			}
			body.WriteString("=== " + res.Name + " ===\n")
		}
		body.WriteString(res.Art)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.txt"`)
	w.Header().Set("Content-Length", strconv.Itoa(body.Len()))
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(body.String())); err != nil {
		log.Println("download write error:", err)
	}
}
