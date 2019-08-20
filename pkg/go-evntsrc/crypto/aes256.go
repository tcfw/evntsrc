package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

//NewAES256 creates a new encrypter/decrypter using
//	Signature: NONE
//	Encryption: AES 256 GCM
func NewAES256(key []byte) (*AES256, error) {
	return &AES256{aesKey: key}, nil
}

//AES256 provides only encryption. NO signatures will be added to events
type AES256 struct {
	aesKey []byte
}

//Encrypt encrypts a payload and provides any additional metadata such as iv, nonce and/or signatures
func (ecl *AES256) Encrypt(in []byte) ([]byte, map[string]string, error) {
	md := map[string]string{}

	//Encrypt payload
	cipherData, nonce, err := ecl.aesEncrypt(in)
	if err != nil {
		return nil, md, err
	}
	md[mdAESNonce] = hex.EncodeToString(nonce)

	return cipherData, md, nil
}

func (ecl *AES256) aesEncrypt(in []byte) ([]byte, []byte, error) {
	aesBlock, err := aes.NewCipher(ecl.aesKey)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	aesgcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, nil, err
	}
	cipherdata := aesgcm.Seal(nil, nonce, in, nil)

	return cipherdata, nonce, nil
}

//Decrypt decrypts the payload given metadata i.e. nonce/iv
func (ecl *AES256) Decrypt(in []byte, md map[string]string) ([]byte, error) {
	aesBlock, err := aes.NewCipher(ecl.aesKey)
	if err != nil {
		return nil, err
	}

	nonceHex, ok := md[mdAESNonce]
	if !ok {
		return nil, fmt.Errorf("MD has no nonce %s", mdAESNonce)
	}
	nonce, err := hex.DecodeString(nonceHex)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, in, nil)
}

//Verify validates the payload given metadata
//Will only check the existance of the AES nonce
func (ecl *AES256) Verify(in []byte, md map[string]string) error {
	if _, ok := md[mdAESNonce]; !ok {
		return fmt.Errorf("No nonce found")
	}
	return nil
}
