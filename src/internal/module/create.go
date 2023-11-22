package module

import (
	"path/filepath"
	"strings"

	"registry-stable/internal/files"
)

func CreateMetadataFile(m Module, moduleDataDir string) error {
	repositoryFileData, err := BuildMetadataFile(m)
	if err != nil {
		return err
	}

	filePath := getFilePath(m, moduleDataDir)
	return files.SafeWriteObjectToJsonFile(filePath, repositoryFileData)
}

func getFilePath(m Module, moduleDataDir string) string {
	return filepath.Join(moduleDataDir, strings.ToLower(m.Namespace[0:1]), m.Namespace, m.Name, m.TargetSystem+".json")
}
