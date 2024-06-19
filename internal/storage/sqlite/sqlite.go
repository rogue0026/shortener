package sqlite

import (
	"database/sql"
	"fmt"

	// _ "github.com/mattn/go-sqlite3"
	"github.com/rogue0026/shortener/internal/config"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Shortener) (*Storage, error) {
	const fn = "storage.sqlite.New"

	connPool, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	tx, err := connPool.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: begin transaction: %w", fn, err)
	}
	query := `CREATE TABLE IF NOT EXISTS urls (
    id INTEGER PRIMARY KEY,
    long_url TEXT NOT NULL,
    short_url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_id ON urls (id);
	CREATE TABLE IF NOT EXISTS users (
		user_id CHAR(23) PRIMARY KEY,
		login VARCHAR(25) UNIQUE,
		password VARCHAR(100),
		email VARCHAR(50) UNIQUE);
	CREATE INDEX IF NOT EXISTS idx_id ON users (user_id);`
	if _, err := tx.Exec(query); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: exec: %w", fn, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: commit: %w", fn, err)
	}

	s := Storage{
		db: connPool,
	}
	return &s, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}
