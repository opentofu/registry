package module

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Version struct {
	Version string `json:"version"` // The version number of the provider. Correlates to a tag in the module repository
}

type MetadataFile struct {
	Versions []Version `json:"versions"`
}

type Module struct {
	Namespace    string // The module namespace
	Name         string // The module name
	TargetSystem string // The module target system
	Directory    string // The root directory that the module lives in
}

func (m Module) RepositoryURL() string {
	return fmt.Sprintf("https://github.com/%s/terraform-%s-%s", m.Namespace, m.TargetSystem, m.Name)
}

// the file should just contain a link to GitHub to download the tarball, ie:
// git::https://github.com/terraform-aws-modules/terraform-aws-iam?ref=v5.30.0
func (m Module) VersionDownloadURL(version Version) string {
	return fmt.Sprintf("git::%s?ref=%s", m.RepositoryURL(), version.Version)
}

func (m Module) MetadataPath() string {
	return filepath.Join(m.Directory, m.Namespace[0:1], m.Namespace, m.Name, m.TargetSystem+".json")
}

func (m Module) ReadMetadata() (MetadataFile, error) {
	var metadata MetadataFile

	path := m.MetadataPath()

	metadataFile, err := os.ReadFile(path)
	if err != nil {
		return metadata, fmt.Errorf("failed to open metadata file: %w", err)
	}

	err = json.Unmarshal(metadataFile, &metadata)
	if err != nil {
		return metadata, fmt.Errorf("failed to unmarshal metadata file: %w", err)
	}

	return metadata, nil
}
