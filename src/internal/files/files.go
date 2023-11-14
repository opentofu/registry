package files

import (
	"encoding/json"
	"io/fs"
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

// TODO: Remove the WriteToFile usages above and replace with this abstraction

type FileWriter interface {
	// WriteFile writes data to a file named by filename.
	WriteFile(name string, data []byte, perm os.FileMode) error

	// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
	MkdirAll(path string, perm os.FileMode) error
}

type RealFileSystem struct {
	fs.FS
}

func (fs *RealFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (fs *RealFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
