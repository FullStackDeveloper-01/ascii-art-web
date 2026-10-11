package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestASCIIArtHandlerRejectsGET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art", nil)
	res := httptest.NewRecorder()

	ASCIIArtHandler(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusMethodNotAllowed)
	}
	if got := res.Header().Get("Allow"); got != http.MethodPost {
		t.Fatalf("Allow = %q, want %q", got, http.MethodPost)
	}
	if !strings.Contains(res.Body.String(), "405 Method Not Allowed") {
		t.Fatalf("response does not contain method error: %q", res.Body.String())
	}
}

func TestASCIIArtHandlerRejectsTextOver100Characters(t *testing.T) {
	form := strings.NewReader("text=" + strings.Repeat("a", maxLength+1) + "&banner=standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()

	ASCIIArtHandler(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
	}
	if !strings.Contains(res.Body.String(), "the limit is 100") {
		t.Fatalf("response does not contain length limit: %q", res.Body.String())
	}
}
