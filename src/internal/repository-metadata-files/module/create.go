package module

import (
	"path/filepath"
	"registry-stable/internal/files"
	"registry-stable/internal/module"
	"strings"
)

func CreateMetadataFile(m module.Module, moduleDataDir string) error {
	repositoryFileData, err := BuildMetadataFile(m)
	if err != nil {
		return err
	}

	filePath := getFilePath(m, moduleDataDir)
	return files.WriteToFile(filePath, repositoryFileData)
}

func getFilePath(m module.Module, moduleDataDir string) string {
	return filepath.Join(moduleDataDir, strings.ToLower(m.Namespace[0:1]), m.Namespace, m.Name, m.TargetSystem+".json")
}
