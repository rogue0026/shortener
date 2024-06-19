package random

import (
	"math/rand"
	"time"
)

var runes = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func String() string {
	seed := rand.NewSource(time.Now().UnixNano())
	g := rand.New(seed)
	res := make([]rune, 8)
	for i := range res {
		res[i] = runes[g.Intn(len(runes))]
	}
	return string(res)
}
