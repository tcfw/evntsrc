package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestECDSASHA384ECDHAES256(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	remoteKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	x509cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "evntsrc.io",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * time.Second),
		PublicKeyAlgorithm:    x509.ECDSA,
		PublicKey:             &remoteKey.PublicKey,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
		BasicConstraintsValid: true,
	}

	crypter := &ECDSASHA384ECDHAES256{
		ECDSASHA384AES256: &ECDSASHA384AES256{
			privKey: privateKey,
		},
		pubCert: x509cert,
	}

	err = crypter.ecdh()
	assert.NoError(t, err)
	assert.NotEmpty(t, crypter.aesKey)

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

func TestMismatchCurves(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	remoteKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	x509cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "evntsrc.io",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * time.Second),
		PublicKeyAlgorithm:    x509.ECDSA,
		PublicKey:             &remoteKey.PublicKey,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
		BasicConstraintsValid: true,
	}

	crypter := &ECDSASHA384ECDHAES256{
		ECDSASHA384AES256: &ECDSASHA384AES256{
			privKey: privateKey,
		},
		pubCert: x509cert,
	}

	err = crypter.ecdh()
	assert.Error(t, err)
}

func TestCertValidation(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	remoteKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	x509template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "evntsrc.io",
		},
		NotBefore:             time.Now().Add(-1 * time.Second),
		NotAfter:              time.Now().Add(1 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, x509template, x509template, &remoteKey.PublicKey, remoteKey)
	if err != nil {
		t.Fatal(err)
	}

	x509cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}

	crypter := &ECDSASHA384ECDHAES256{
		ECDSASHA384AES256: &ECDSASHA384AES256{
			privKey: privateKey,
		},
		pubCert: x509cert,
	}

	err = crypter.validatepubCert()
	assert.Error(t, err)
	//Assert self signed for now
	assert.IsType(t, x509.UnknownAuthorityError{}, err)
}

func TestFetchpubCert(t *testing.T) {
	remoteKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	x509template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "evntsrc.io",
		},
		NotBefore:             time.Now().Add(-1 * time.Second),
		NotAfter:              time.Now().Add(1 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, x509template, x509template, &remoteKey.PublicKey, remoteKey)
	if err != nil {
		t.Fatal(err)
	}

	der := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(der)
	}))
	defer ts.Close()

	crypter := &ECDSASHA384ECDHAES256{}

	err = crypter.fetchpubCert(ts.URL)
	assert.Error(t, err)
	//Assert self signed for now
	assert.IsType(t, x509.UnknownAuthorityError{}, err)
}
