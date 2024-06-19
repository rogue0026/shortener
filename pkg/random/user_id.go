package random

import (
	"crypto/rand"
	"math/big"
)

const chars = "abcdefghijklmnopqrstuvwxy0123456789"

func UserID() (string, error) {
	uuid := make([]rune, 11)

	for i := 0; i < len(uuid); i++ {
		bi, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		charPos := int(bi.Int64())
		if err != nil {
			return "", err
		}
		uuid[i] = rune(chars[charPos])
		if i == 5 {
			uuid[i] = '-'
		}
	}

	return string(uuid), nil
}
