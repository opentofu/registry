package files

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func WriteToJsonFile(filePath string, data interface{}) error {
	marshalledJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	err = os.MkdirAll(path.Dir(filePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err = os.WriteFile(filePath, marshalledJson, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
