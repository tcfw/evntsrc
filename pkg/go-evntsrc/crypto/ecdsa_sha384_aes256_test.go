package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestECDSASHA384AES256(t *testing.T) {
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		t.Fatal(err)
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	crypter := &ECDSASHA384AES256{
		aesKey:  aesKey,
		privKey: privateKey,
	}

	msg := []byte("test Data")

	cipher, md, err := crypter.Encrypt(msg)
	assert.NoError(t, err)
	assert.Len(t, md, 4) //nounce, sig-r, sig-n, hash
	assert.NotEqual(t, msg, cipher)

	err = crypter.Verify(cipher, md)
	assert.NoError(t, err)

	restore, err := crypter.Decrypt(cipher, md)
	assert.NoError(t, err)
	assert.Equal(t, restore, msg)
}
