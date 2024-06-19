package login

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rogue0026/shortener/internal/auth/token"
	"github.com/rogue0026/shortener/internal/domain"
	"github.com/rogue0026/shortener/internal/storage"
	"github.com/sirupsen/logrus"
)

type Loginer interface {
	LoginUser(ctx context.Context, inLogin, inPassword string) (string, error)
}

func New(logger *logrus.Logger, loginer Loginer) http.HandlerFunc {
	const fn = "handlers.users.login.New"
	return func(w http.ResponseWriter, r *http.Request) {
		in := domain.User{}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userID, err := loginer.LoginUser(r.Context(), in.Login, in.Password)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrUserNotFound):
				w.WriteHeader(http.StatusNotFound)
				return
			case errors.Is(err, storage.ErrInvalidPassword):
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid login or password"))
				return
			}
		}
		tokenString, err := token.Generate(userID)
		if err != nil {
			logger.Errorf("%s: %s", fn, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Token", tokenString)
		w.WriteHeader(http.StatusOK)
	}
}
