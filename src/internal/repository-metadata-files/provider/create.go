package provider

import (
	"path/filepath"
	"registry-stable/internal/files"
	"registry-stable/internal/provider"
	"strings"
)

func CreateMetadataFile(p provider.Provider, providerDataDir string) error {
	repositoryFileData, err := BuildMetadataFile(p, providerDataDir)
	if err != nil {
		return err
	}

	filePath := getFilePath(p, providerDataDir)
	return files.WriteToFile(filePath, repositoryFileData)
}

func getFilePath(p provider.Provider, providerDataDir string) string {
	return filepath.Join(providerDataDir, strings.ToLower(p.Namespace[0:1]), p.Namespace, p.ProviderName+".json")
}
