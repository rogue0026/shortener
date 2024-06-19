package ping

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type Pinger interface {
	Ping() error
}

func New(p Pinger, logger *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "http-server.handlers.New"
		err := p.Ping()
		if err != nil {
			logger.Errorf("%s: %w", fn, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
