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
	if r.URL.Scheme != "http" {
		r.URL.Scheme = "http"
		r.URL.Host = r.Host
		// badRequest(w, fmt.Errorf("URL scheme is not http"))
	}

	validShortURL := regexp.MustCompile(`^/?[a-zA-Z0-9]+$`)
	shortURL := strings.TrimPrefix(r.URL.Path, "/")

	if !validShortURL.MatchString(shortURL) {
		badRequest(w, fmt.Errorf("invalid shortUrl"))
		return
	}

	longURL, err := h.service.ReStore(shortURL)
	if err != nil {
		if err.Error() == "key not found" {
			badRequest(w, err)
		} else {
			internalError(w, err)
		}
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
	// w.Header().Add("Location", longURL)
	// w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *LSHandlers) HandlePost(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if r.URL.Path != "/" || r.Body == http.NoBody || err != nil || !validURL(string(b)) {
		badRequest(w, fmt.Errorf("bad request"))
		return
	}

	longURL := string(b)
	shortURL, err := h.service.Store(longURL)

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
	_, err = w.Write([]byte(shortURL))
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
