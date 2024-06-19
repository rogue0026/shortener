package register

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rogue0026/shortener/internal/domain"
	"github.com/sirupsen/logrus"
)

type UserRegister interface {
	RegisterUser(ctx context.Context, login, password, email string) (string, error)
}

func New(logger *logrus.Logger, register UserRegister) http.HandlerFunc {
	const fn = "handlers.users.register.New"
	return func(w http.ResponseWriter, r *http.Request) {
		in := domain.User{}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		uuid, err := register.RegisterUser(r.Context(), in.Login, in.Password, in.Email)
		if err != nil {
			logger.Errorf("%s: %s", fn, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(uuid))
	}
}
