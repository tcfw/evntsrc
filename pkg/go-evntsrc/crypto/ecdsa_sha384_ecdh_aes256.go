package crypto

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pubCertURL = "https://evntsrc.io/.well-known/keys/sig.pem"
)

//ECDSASHA384ECDHAES256 provides both encryption and signatures into events
type ECDSASHA384ECDHAES256 struct {
	*ECDSASHA384AES256
	pubCert *x509.Certificate
}

//NewECDSASHA384ECDHAES256 creates a new encrypter/decrypter using
//	Key Generation: ECDH with well known key
//	Signature: ECDSA on a SHA384 hash
//	Encryption: AES 256 GCM
func NewECDSASHA384ECDHAES256(privKey *ecdsa.PrivateKey) (*ECDSASHA384ECDHAES256, error) {
	crypter := &ECDSASHA384ECDHAES256{
		ECDSASHA384AES256: &ECDSASHA384AES256{
			privKey: privKey,
		},
		pubCert: nil,
	}

	err := crypter.fetchpubCert(pubCertURL)
	if err != nil {
		return nil, err
	}

	return crypter, nil
}

func (ecl *ECDSASHA384ECDHAES256) ecdh() error {
	if ecl.pubCert.PublicKey == nil {
		return fmt.Errorf("No pub key")
	}

	if ecl.pubCert.PublicKeyAlgorithm != x509.ECDSA {
		return fmt.Errorf("Unsupported x509 public key algo")
	}

	publicECC := ecl.pubCert.PublicKey.(*ecdsa.PublicKey)

	if !ecl.privKey.IsOnCurve(publicECC.X, publicECC.Y) {
		return fmt.Errorf("Public cert curve does not match private key curve")
	}
	if !publicECC.IsOnCurve(ecl.privKey.X, ecl.privKey.Y) {
		return fmt.Errorf("Private key curve does not match public cert curve")
	}

	skX, _ := publicECC.Curve.ScalarMult(publicECC.X, publicECC.Y, ecl.privKey.D.Bytes())
	key := sha256.Sum256(skX.Bytes())
	ecl.aesKey = key[:]

	return nil
}

func (ecl *ECDSASHA384ECDHAES256) fetchpubCert(url string) error {
	tr := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     10 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	pemBlock, _ := pem.Decode(body)
	if pemBlock == nil {
		return fmt.Errorf("failed to parse PEM cert")
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return err
	}
	ecl.pubCert = cert

	return ecl.validatepubCert()
}

func (ecl *ECDSASHA384ECDHAES256) validatepubCert() error {
	_, err := ecl.pubCert.Verify(x509.VerifyOptions{
		DNSName: "evntsrc.io",
	})
	if err != nil {
		return err
	}

	if ecl.pubCert.PublicKeyAlgorithm != x509.ECDSA {
		return fmt.Errorf("Unsupported x509 public key algo")
	}

	return nil
}
