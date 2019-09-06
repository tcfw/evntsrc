package evntsrc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	esCrypto "github.com/tcfw/evntsrc/pkg/go-evntsrc/crypto"
)

//EncrypterDecrypter encrypts, decrypts and verifies message payloads
type EncrypterDecrypter interface {
	//Encrypt encrypts a payload and provides any additional metadata such as iv, nonce and/or signatures
	Encrypt([]byte) ([]byte, map[string]string, error)

	//Decrypt decrypts the payload given metadata i.e. nonce/iv
	Decrypt([]byte, map[string]string) ([]byte, error)

	//Verify validates the payload given metadata i.e. signatures
	Verify([]byte, map[string]string) error
}

//EphemeralCrypto generates ephemeral crypto for use in testing
func EphemeralCrypto() (*esCrypto.ECDSASHA384AES256, error) {
	ephemAESKey := make([]byte, 32)
	if _, err := rand.Read(ephemAESKey); err != nil {
		return nil, err
	}

	ephemPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return esCrypto.NewECDSASHA384AES256(ephemAESKey, ephemPriv)
}
