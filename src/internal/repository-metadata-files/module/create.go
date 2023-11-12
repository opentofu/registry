package module

import (
	"fmt"
	"registry-stable/internal/files"
	"registry-stable/internal/module"
)

func CreateMetadataFile(m module.Module) error {
	repositoryFileData, err := BuildMetadataFile(m)
	if err != nil {
		return err
	}

	filePath := getFilePath(m)
	return files.WriteToFile(filePath, repositoryFileData)
}

func getFilePath(m module.Module) string {
	shard := m.Namespace[0]
	return fmt.Sprintf("modules/%c/%s/%s/%s.json", shard, m.Namespace, m.Name, m.System)
}
