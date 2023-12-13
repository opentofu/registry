package base

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/opentofu/registry-stable/internal/files"
)

type FileSystem struct {
	Directory string
}

func (fs FileSystem) ReadJSONFrom(path string, ref any) error {
	data, err := os.ReadFile(filepath.Join(fs.Directory, path))
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, ref)
	if err != nil {
		return err
	}
	return nil
}

func (fs FileSystem) WriteJSONInto(path string, data any) error {
	return files.SafeWriteObjectToJSONFile(filepath.Join(fs.Directory, path), data)
}

func (fs FileSystem) List() ([]string, error) {
	paths := make([]string, 0)
	err := filepath.Walk(fs.Directory, func(path string, info os.FileInfo, err error) error {
		paths = append(paths, path)
		return nil
	})
	return paths, err
}
