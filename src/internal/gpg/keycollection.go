package gpg

import (
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
	location := k.MetadataPath()

	// check if the directory exists
	if _, err := os.Stat(location); os.IsNotExist(err) {
		return nil, nil
	}

	// if it does exist, iterate across the files
	files, err := os.ReadDir(location)
	if err != nil {
		return nil, err
	}

	keys := make([]Key, len(files))
	for i, file := range files {
		key, err := buildKey(filepath.Join(location, file.Name()))
		if err != nil {
			return nil, err
		}
		keys[i] = *key
	}

	return keys, nil
}
