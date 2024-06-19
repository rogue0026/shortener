package token

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rogue0026/shortener/internal/auth"
)

type Claims struct {
	UserID string
	jwt.RegisteredClaims
}

func Generate(userID string) (string, error) {
	c := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	key := os.Getenv("PRIVATE_KEY")
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func Verify(tokenString string) error {
	key := os.Getenv("PRIVATE_KEY")

	c := Claims{}
	_, err := jwt.ParseWithClaims(tokenString, &c, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, auth.ErrInvalidSigningMethod
		}
		return []byte(key), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return auth.ErrTokenIsExpired
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return auth.ErrTokenIsInvalid
		}
	}

	return nil
}
