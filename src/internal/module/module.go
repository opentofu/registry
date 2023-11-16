package module

import (
	"fmt"
	"path/filepath"
	"registry-stable/internal"
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
}

func (m Module) RepositoryURL() string {
	return fmt.Sprintf("https://github.com/%s/terraform-%s-%s", m.Namespace, m.TargetSystem, m.Name)
}

// the file should just contain a link to GitHub to download the tarball, ie:
// git::https://github.com/terraform-aws-modules/terraform-aws-iam?ref=v5.30.0
func (m Module) VersionDownloadURL(version Version) string {
	return fmt.Sprintf("git::%s?ref=%s", m.RepositoryURL(), version.Version)
}

func (m Module) MetadataPath(directory string) string {
	return filepath.Join(directory, m.Namespace[0:1], m.Namespace, m.Name, m.TargetSystem+".json")
}

func (m Module) outputPath(directory string) string {
	return filepath.Join(directory, "v1", "modules", m.Namespace, m.Name, m.TargetSystem)
}

func (m Module) VersionListingPath(directory string) string {
	return filepath.Join(m.outputPath(directory), "versions")
}

func (m Module) VersionDownloadPath(directory string, v Version) string {
	return filepath.Join(m.outputPath(directory), internal.TrimTagPrefix(v.Version), "download")
}
