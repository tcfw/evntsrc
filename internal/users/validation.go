package users

import (
	"crypto/rand"
)

func genValidationToken() string {
	n := 64
	runes := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-.:`

	rb := make([]byte, n)
	rand.Read(rb)

	b := make([]byte, n)
	for i, by := range rb {
		b[i] = runes[by%byte(len(runes))]
	}
	return string(b)
}
