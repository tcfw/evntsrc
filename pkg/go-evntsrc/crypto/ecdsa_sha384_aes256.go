package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
)

//NewECDSASHA384AES256 creates a new encrypter/decrypter using
//	Signature: ECDSA on a SHA384 hash
//	Encryption: AES 256 GCM
func NewECDSASHA384AES256(key []byte, privKey *ecdsa.PrivateKey) (*ECDSASHA384AES256, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("Key too short. Key length must be 32 bytes")
	}

	return &ECDSASHA384AES256{
		aesKey:  key,
		privKey: privKey,
	}, nil
}

//ECDSASHA384AES256 provides both encryption and signatures into events
type ECDSASHA384AES256 struct {
	aesKey  []byte
	privKey *ecdsa.PrivateKey
}

//Encrypt transforms in bytes into an ECDSA signature with an AES 256 cipher
func (ecl *ECDSASHA384AES256) Encrypt(in []byte) ([]byte, map[string]string, error) {
	md := map[string]string{}

	//Encrypt payload
	cipherData, nonce, err := ecl.aesEncrypt(in)
	if err != nil {
		return nil, md, err
	}
	md[mdAESNonce] = hex.EncodeToString(nonce)

	//Hash
	hash := sha512.Sum384(in)
	md[mdHash] = hex.EncodeToString(hash[:])

	//Sign
	r, s, err := ecdsa.Sign(rand.Reader, ecl.privKey, hash[:])
	if err != nil {
		return nil, md, err
	}
	md[mdSigR] = r.Text(16)
	md[mdSigS] = s.Text(16)

	return cipherData, md, nil
}

func (ecl *ECDSASHA384AES256) aesEncrypt(in []byte) ([]byte, []byte, error) {
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

//Decrypt decrypts the AES 256 payload given metadata of nonce
func (ecl *ECDSASHA384AES256) Decrypt(in []byte, md map[string]string) ([]byte, error) {
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

//Verify validates the payload given metadata sig-r and sig-s
func (ecl *ECDSASHA384AES256) Verify(in []byte, md map[string]string) error {
	rHex, ok := md[mdSigR]
	if !ok {
		return fmt.Errorf("MD missing signature part %s", mdSigR)
	}

	sHex, ok := md[mdSigS]
	if !ok {
		return fmt.Errorf("MD missing signature part %s", mdSigS)
	}

	hashHex, ok := md[mdHash]
	if !ok {
		return fmt.Errorf("MD missing signature part %s", mdHash)
	}

	hash, err := hex.DecodeString(hashHex)
	if err != nil {
		return err
	}

	r := new(big.Int)
	s := new(big.Int)

	r.SetString(rHex, 16)
	s.SetString(sHex, 16)

	valid := ecdsa.Verify(&ecl.privKey.PublicKey, hash, r, s)
	if !valid {
		return fmt.Errorf("ECDSA Signature invalid (%s, %s, %s)", rHex, sHex, hashHex)
	}

	return nil
}
