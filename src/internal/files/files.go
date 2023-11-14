package files

import (
	"encoding/json"
	"os"
	"path"
)

func WriteToFile(filePath string, data interface{}) error {
	marshalledJson, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, marshalledJson, 0644)
	if err != nil {
		return err
	}

	return nil
}
