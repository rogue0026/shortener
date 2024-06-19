package get

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rogue0026/shortener/internal/storage"
	"github.com/sirupsen/logrus"
)

type URLFetcher interface {
	FetchLongURL(shortURL string) (string, error)
}

func New(logger *logrus.Logger, fetcher URLFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "short_url")
		longURL, err := fetcher.FetchLongURL(shortURL)
		if err != nil {
			if errors.Is(err, storage.ErrRowNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				logger.Errorln(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

		return
	}
}
