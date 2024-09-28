package web

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/patraden/ya-practicum-go-shortly/internal/services"
)

func badRequest(w http.ResponseWriter, e error) {
	http.Error(w, e.Error(), http.StatusBadRequest)
}

func internalError(w http.ResponseWriter, e error) {
	http.Error(w, e.Error(), http.StatusInternalServerError)
}

type LSHandlers struct {
	service *services.LinkStore
}

func NewLSHandlers(ls *services.LinkStore) *LSHandlers {
	return &LSHandlers{
		service: ls,
	}
}

func (h *LSHandlers) HandleRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.HandlePost(w, r)
	} else if r.Method == http.MethodGet {
		h.HandleGet(w, r)
	} else {
		badRequest(w, fmt.Errorf("unknown method"))
	}
}

func (h *LSHandlers) HandleGet(w http.ResponseWriter, r *http.Request) {
	validShortURL := regexp.MustCompile(`^/?[a-zA-Z0-9]+$`)
	shortUrl := strings.TrimPrefix(r.URL.Path, "/")

	if !validShortURL.MatchString(shortUrl) {
		badRequest(w, fmt.Errorf("invalid shortUrl"))
		return
	}

	longUrl, err := h.service.ReStore(shortUrl)
	if err != nil {
		if err.Error() == "key not found" {
			badRequest(w, err)
		} else {
			internalError(w, err)
		}
		return
	}

	// http.Redirect(w, r, longUrl, http.StatusTemporaryRedirect)
	w.Header().Add("Location", longUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *LSHandlers) HandlePost(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if r.URL.Path != "/" || r.Body == http.NoBody || err != nil || !validURL(string(b)) {
		badRequest(w, fmt.Errorf("bad request"))
		return
	}

	longUrl := string(b)
	shortUrl, err := h.service.Store(longUrl)

	if err != nil {
		if err.Error() == "out of shortlinks" {
			internalError(w, err)
		} else {
			badRequest(w, err)
		}
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortUrl))
	if err != nil {
		internalError(w, err)
		return
	}

}

func validURL(s string) bool {
	parsedURL, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}
	return true
}
