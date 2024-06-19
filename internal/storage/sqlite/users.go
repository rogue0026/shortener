package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rogue0026/shortener/internal/domain"
	"github.com/rogue0026/shortener/internal/storage"
	"github.com/rogue0026/shortener/pkg/random"
	"golang.org/x/crypto/bcrypt"
)

func (s *Storage) RegisterUser(ctx context.Context, login, password, email string) (string, error) {
	const fn = "storage.sqlite.RegisterUser"
	query := `INSERT INTO users (user_id, login, password, email) VALUES (?, ?, ?, ?)`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("%s: prepare: %w", fn, err)
	}
	uuid, err := random.UserID()
	if err != nil {
		return "", fmt.Errorf("%s: generating uuid: %w", fn, err)
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	if _, err := stmt.ExecContext(ctx, uuid, login, string(passHash), email); err != nil {
		return "", fmt.Errorf("%s: exec: %w", fn, err)
	}

	return uuid, nil
}

func (s *Storage) LoginUser(ctx context.Context, inLogin, inPassword string) (string, error) {
	const fn = "storage.sqlite.LoginUser"

	query := `SELECT user_id, login, password FROM users WHERE login = ?;`
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return "", err
	}

	var u = domain.User{}
	if err := stmt.QueryRowContext(ctx, inLogin).Scan(&u.UserID, &u.Login, &u.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUserNotFound
		}
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	if err := u.CheckPassword(inPassword); err != nil {
		return "", storage.ErrInvalidPassword
	}

	return u.UserID, nil
}
