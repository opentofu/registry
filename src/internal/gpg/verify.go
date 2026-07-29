package gpg

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	openpgpErrors "github.com/ProtonMail/go-crypto/openpgp/errors"
)

func VerifyDetachedSignature(keys []Key, data, signature []byte) (bool, error) {
	for _, key := range keys {
		keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(key.ASCIIArmor))
		if err != nil {
			return false, fmt.Errorf("failed to parse key %s: %w", key.KeyID, err)
		}

		_, err = openpgp.CheckDetachedSignature(
			keyring,
			bytes.NewReader(data),
			bytes.NewReader(signature),
			nil,
		)

		if err == nil {
			return true, nil
		}

		// This key didn't produce the signature; try the next one.
		if errors.Is(err, openpgpErrors.ErrUnknownIssuer) {
			continue
		}

		// Expired keys are still accepted by the registry.
		if errors.Is(err, openpgpErrors.ErrKeyExpired) || errors.Is(err, openpgpErrors.ErrSignatureExpired) {
			return true, nil
		}

		return false, fmt.Errorf("error checking signature with key %s: %w", key.KeyID, err)
	}

	return false, nil
}
