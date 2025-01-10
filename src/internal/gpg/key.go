package gpg

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

type Key struct {
	ASCIIArmor string `json:"ascii_armor"`
	KeyID      string `json:"key_id"`
}

func buildKey(path string) (*Key, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open key file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read key file: %w", err)
	}

	asciiArmor := string(data)

	key, err := ParseKey(asciiArmor)
	if err != nil {
		return nil, fmt.Errorf("could not parse key: %w", err)
	}

	return &Key{
		ASCIIArmor: asciiArmor,
		KeyID:      strings.ToUpper(key.GetHexKeyID()),
	}, nil
}

// ParseKey parses a GPG key from ascii armor.
func ParseKey(data string) (*crypto.Key, error) {
	key, err := crypto.NewKeyFromArmored(data)
	if err != nil {
		return nil, fmt.Errorf("could not build public key from ascii armor: %w", err)
	}

	return key, nil
}
