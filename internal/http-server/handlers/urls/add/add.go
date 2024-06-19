package add

import (
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type URLInserter interface {
	InsertShortURL(longURL string) (string, error)
}

func New(logger *logrus.Logger, inserter URLInserter, address string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Errorln(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		longURL := string(data)
		shortURL, err := inserter.InsertShortURL(longURL)
		if err != nil {
			logger.Errorln(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("%s/%s\n", address, shortURL)))
	}
}
