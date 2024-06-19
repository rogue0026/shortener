package auth

import "errors"

var (
	ErrTokenIsExpired       = errors.New("token is expired")
	ErrTokenIsInvalid       = errors.New("token is invalid")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)
