package provider

import (
	"path/filepath"
	"strings"

	"registry-stable/internal/files"
	"registry-stable/internal/provider"
)

func CreateMetadataFile(p provider.Provider, providerDataDir string) error {
	repositoryFileData, err := BuildMetadataFile(p, providerDataDir)
	if err != nil {
		return err
	}

	filePath := getFilePath(p, providerDataDir)
	return files.SafeWriteObjectToJsonFile(filePath, repositoryFileData)
}

func getFilePath(p provider.Provider, providerDataDir string) string {
	return filepath.Join(providerDataDir, strings.ToLower(p.Namespace[0:1]), p.Namespace, p.ProviderName+".json")
}
