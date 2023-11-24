package files

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// SafeWriteObjectToJSONFile writes the given object to the given file path.
// It also ensures that the destination directory exists and that the file is written correctly.
func SafeWriteObjectToJSONFile(filePath string, data interface{}) error {
	marshalledJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal for %s: %w", filePath, err)
	}

	err = os.MkdirAll(path.Dir(filePath), 0755) // nolint: gomnd // 0755 is the default for os.MkdirAll
	if err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", filePath, err)
	}

	err = os.WriteFile(filePath, marshalledJSON, 0600)
	if err != nil {
		// Error already contains filePath
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
