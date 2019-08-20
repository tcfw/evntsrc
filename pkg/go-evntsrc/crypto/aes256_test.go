package crypto

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES256(t *testing.T) {
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		t.Fatal(err)
	}

	crypter := &AES256{
		aesKey: aesKey,
	}

	msg := []byte("test Data")

	cipher, md, err := crypter.Encrypt(msg)
	assert.NoError(t, err)
	assert.Len(t, md, 1) //nounce
	assert.NotEqual(t, msg, cipher)

	err = crypter.Verify(cipher, md)
	assert.NoError(t, err)

	restore, err := crypter.Decrypt(cipher, md)
	assert.NoError(t, err)
	assert.Equal(t, restore, msg)
}
