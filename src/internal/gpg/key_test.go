package gpg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/stretchr/testify/assert"
)

func generatePrivateKey() (string, error) {
	// Generate a new RSA private key with 2048 bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}
	// Encode the private key to the PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	return string(pem.EncodeToMemory(privateKeyPEM)), nil
}

func generateGPGKey() (string, error) {
	rsaKey, err := helper.GenerateKey("test", "test", []byte("test"), "rsa", 1024)
	if err != nil {
		panic(err)
	}

	keyRing, err := crypto.NewKeyFromArmoredReader(strings.NewReader(rsaKey))
	if err != nil {
		panic(err)
	}

	publicKey, err := keyRing.GetArmoredPublicKey()
	if err != nil {
		panic(err)
	}

	return publicKey, nil
}

func TestParseKey(t *testing.T) {
	stringPtr := func(s string) *string {
		return &s
	}

	privateKey, _ := generatePrivateKey()
	publicGPGKey, _ := generateGPGKey()

	tests := []struct {
		name        string
		data        string
		expectedErr *string
	}{
		{
			name:        "public gpg key should succeed",
			data:        publicGPGKey,
			expectedErr: nil,
		},
		{
			name:        "private key should fail",
			data:        privateKey,
			expectedErr: stringPtr("could not build public key from ascii armor"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseKey(test.data)

			if test.expectedErr != nil {
				assert.ErrorContains(t, err, *test.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
