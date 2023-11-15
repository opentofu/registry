package module

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
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

func (m Module) Logger(logger *slog.Logger) *slog.Logger {
	return logger.With(slog.String("namespace", m.Namespace), slog.String("name", m.Name), slog.String("targetSystem", m.TargetSystem))
}

func (m Module) repositoryPath() string {
	return fmt.Sprintf("%s/terraform-%s-%s", m.Namespace, m.TargetSystem, m.Name)
}

func (m Module) RepositoryURL() string {
	return fmt.Sprintf("https://github.com/%s", m.repositoryPath())
}

func (m Module) MetadataPath(directory string) string {
	return filepath.Join(directory, m.Namespace[0:1], m.Namespace, m.Name, m.TargetSystem+".json")
}

// ReadMetadata reads the module metadata file from the filesystem directly. This data should be the data fetched from the git repository.
func (m Module) ReadMetadata(moduleDirectory string, logger *slog.Logger) (*MetadataFile, error) {
	path := m.MetadataPath(moduleDirectory)

	// open the file
	metadataFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}

	// Read the file contents into a Module[] struct
	var metadata MetadataFile
	err = json.Unmarshal(metadataFile, &metadata)
	if err != nil {
		return nil, err
	}

	logger.Debug("Loaded Module Versions", slog.Any("count", len(metadata.Versions)))

	return &metadata, nil
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

// the file should just contain a link to GitHub to download the tarball, ie:
// git::https://github.com/terraform-aws-modules/terraform-aws-iam?ref=v5.30.0
func (m Module) VersionDownloadURL(version Version) string {
	return fmt.Sprintf("git::github.com/%s?ref=%s", m.repositoryPath(), version.Version)
}
