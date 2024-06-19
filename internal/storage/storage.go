package storage

import "errors"

var (
	ErrRowNotFound     = errors.New("requested row not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("requested user not found")
)
