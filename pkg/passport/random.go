package passport

import (
	"crypto/rand"
	"encoding/hex"
)

//RandomString generates a n lengthed string (cryptographically)
func RandomString(n int) string {
	var randomBytes = make([]byte, n/2)
	rand.Read(randomBytes)

	return hex.EncodeToString(randomBytes)
}
