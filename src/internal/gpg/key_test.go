package gpg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func generateExpiredGPGKey(t *testing.T) *crypto.Key {
	t.Helper()

	config := &packet.Config{
		Time:            func() time.Time { return time.Now().Add(-48 * time.Hour) },
		KeyLifetimeSecs: 60,
	}

	entity, err := openpgp.NewEntity("test", "expired test key", "test@test", config)
	require.NoError(t, err)

	key, err := crypto.NewKeyFromEntity(entity)
	require.NoError(t, err)

	return key
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

func TestCanSign(t *testing.T) {
	t.Run("valid key can sign and verify", func(t *testing.T) {
		armored, err := generateGPGKey()
		require.NoError(t, err)

		key, err := ParseKey(armored)
		require.NoError(t, err)

		assert.False(t, key.IsExpired())
		assert.True(t, key.CanVerify(), "sanity: non-expired key should also pass CanVerify")
		assert.True(t, CanSign(key))
	})

	t.Run("expired key can still sign", func(t *testing.T) {
		key := generateExpiredGPGKey(t)

		require.True(t, key.IsExpired(), "test fixture should be expired")
		require.False(t, key.CanVerify(), "CanVerify rejects expired keys, which is what CanSign works around")

		assert.True(t, CanSign(key), "expired keys are still accepted by the registry")
	})
}
