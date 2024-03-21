package gpg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// KeyCollection represents the GPG keys stored in the registry for a specific namespace.
type KeyCollection struct {
	Namespace string // The key namespace
	Directory string // The root directory that the key lives in
}

func (k KeyCollection) MetadataPath() string {
	firstChar := strings.ToLower(k.Namespace[0:1])
	return filepath.Join(k.Directory, firstChar, k.Namespace)
}

func (k KeyCollection) ListKeys() ([]Key, error) {
	location := strings.ToLower(k.MetadataPath())

	// check if the directory exists
	if _, err := os.Stat(location); os.IsNotExist(err) {
		return nil, nil
	}

	// if it does exist, iterate across the files
	files, err := os.ReadDir(location)
	if err != nil {
		return nil, fmt.Errorf("error reading directory %s: %w", location, err)
	}

	keys := make([]Key, len(files))
	for i, file := range files {
		keyPath := filepath.Join(location, file.Name())
		key, err := buildKey(keyPath)
		if err != nil {
			return nil, fmt.Errorf("error building key at %s: %w", keyPath, err)
		}
		keys[i] = *key
	}

	return keys, nil
}
