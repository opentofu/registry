package provider

import (
	"fmt"
	"registry-stable/internal/files"
	"registry-stable/internal/provider"
)

func CreateMetadataFile(p provider.Provider) error {
	repositoryFileData, err := BuildMetadataFile(p)
	if err != nil {
		return err
	}

	filePath := getFilePath(p)
	return files.WriteToFile(filePath, repositoryFileData)
}

func getFilePath(p provider.Provider) string {
	shard := p.Namespace[0]
	return fmt.Sprintf("providers/%c/%s/%s.json", shard, p.Namespace, p.ProviderName)
}
