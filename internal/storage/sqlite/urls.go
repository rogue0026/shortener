package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/rogue0026/shortener/internal/storage"
	"github.com/rogue0026/shortener/pkg/random"
)

func (s *Storage) InsertShortURL(longURL string) (string, error) {
	const fn = "storage.sqlite.InsertShortURL"

	query := `INSERT INTO urls (long_url, short_url) VALUES (?, ?);`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s: prepare: %w", fn, err)
	}

	shortURL := random.String()
	if _, err := stmt.Exec(longURL, shortURL); err != nil {
		return "", fmt.Errorf("%s: exec: %w", fn, err)
	}

	return shortURL, nil
}

func (s *Storage) FetchLongURL(shortURL string) (string, error) {
	const fn = "storage.sqlite.FetchLongURL"

	query := `SELECT long_url FROM urls WHERE short_url = ?;`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("%s: prepare: %w", fn, err)
	}

	var longURL string
	if err := stmt.QueryRow(shortURL).Scan(&longURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrRowNotFound
		} else {
			return "", fmt.Errorf("%s: query: %w", fn, err)
		}
	}

	return longURL, nil
}
