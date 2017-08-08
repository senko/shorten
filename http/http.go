package http

import (
	".."
	"net/http"
)

type Expander struct {
	shortener *short.Shortener
	http.Handler
}

func NewExpander(s *short.Shortener) *Expander {
	return &Expander{
		shortener: s,
	}
}

func (e *Expander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := e.shortener.Expand(r.URL.Path[1:])
	if url == "" {
		http.Error(w, "Short URL not found", http.StatusNotFound)
	} else {
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}
