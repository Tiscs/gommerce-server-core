package secure

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	DEFAULT_PASSWORD_SYMBOLS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+,.?/:;{}[]`~"
)

// RandString randomly generates l length string from the given symbols.
func RandString(l int, s string) (string, error) {
	if s == "" || l <= 0 {
		return "", nil
	}
	c := []rune(s)
	m := big.NewInt(int64(len(c)))
	var r = make([]rune, l)
	for i := 0; i < l; i++ {
		if v, err := rand.Int(rand.Reader, m); err == nil {
			r[i] = c[v.Int64()]
		} else {
			return "", err
		}
	}
	return string(r), nil
}

// MaskString masks the given string.
// The middle third of the given string will be masked with asterisks.
func MaskString(origin string) string {
	l := len(origin)
	r := l % 3
	s := l / 3
	e := s * 2
	if r > 0 {
		e++
	}
	return origin[:s] + strings.Repeat("*", e-s) + origin[e:]
}
