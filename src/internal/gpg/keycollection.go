package gpg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// KeyCollection represents the GPG keys stored in the registry for a specific namespace and provider.
type KeyCollection struct {
	Namespace    string // The key namespace
	ProviderName string // The key provider name
	Directory    string // The root directory that the key lives in
}

func (k KeyCollection) NamespacePath() string {
	firstChar := strings.ToLower(k.Namespace[0:1])
	return filepath.Join(k.Directory, firstChar, k.Namespace)
}

func (k KeyCollection) ProviderPath() string {
	return filepath.Join(k.NamespacePath(), k.ProviderName)
}

func (k KeyCollection) ListKeys() ([]Key, error) {
	namespaceKeys, namespaceErr := k.listKeysIn(k.NamespacePath())
	if namespaceErr != nil {
		return nil, namespaceErr
	}
	providerKeys, providerErr := k.listKeysIn(k.ProviderPath())
	if providerErr != nil {
		return nil, providerErr
	}
	return append(namespaceKeys, providerKeys...), nil
}
func (k KeyCollection) listKeysIn(location string) ([]Key, error) {
	// check if the directory exists
	if _, err := os.Stat(location); os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// if it does exist, iterate across the files
	files, err := os.ReadDir(location)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", location, err)
	}

	keys := make([]Key, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		keyPath := filepath.Join(location, file.Name())
		key, err := buildKey(keyPath)
		if err != nil {
			return nil, fmt.Errorf("error building key at %s: %w", keyPath, err)
		}
		keys = append(keys, *key)
	}

	return keys, nil
}
