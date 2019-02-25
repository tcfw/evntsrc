package passport

import (
	"encoding/hex"
	fmt "fmt"
	"math/rand"
)

func validateSNounce(jti string, snounce string) (string, error) {
	return genNounce()
	// if err != nil {
	// 	return "", err
	// }
	// return "", nil
}

func genNounce() (string, error) {
	var randomBytes = make([]byte, 64)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to create random bytes for nounce")
	}

	return hex.EncodeToString(randomBytes), nil
}
