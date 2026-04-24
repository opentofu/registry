// Package module provides module metadata and build functionality.
package module

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opentofu/registry-stable/internal/blacklist"
	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/github"
)

// Version represents a single version of a module.
type Version struct {
	Version    string     `json:"version"`              // The version number of the provider. Correlates to a tag in the module repository
	Commit     string     `json:"commit"`               // The commit hash of the version tag when the module version was first discovered
	Discovered *time.Time `json:"discovered,omitempty"` // The date the module version was first discovered
}

// Any module without a discovered date should default to the date we started tracking discovery
var defaultDiscovery, _ = time.Parse(time.RFC3339, "2026-04-21T00:00:00Z")

func (v *Version) FirstDiscovered() time.Time {
	if v.Discovered == nil {
		// This can be removed once the backfill is complete
		return defaultDiscovery
	}
	return *v.Discovered
}

// Metadata represents all the metadata for a module. This includes the list of
// versions available for the module.
type Metadata struct {
	Versions []Version `json:"versions"`
}

// Module represents a single module.
type Module struct {
	Namespace    string               // The module namespace
	Name         string               // The module name
	TargetSystem string               // The module target system
	Directory    string               // The root directory that the module lives in
	Logger       *slog.Logger         // A logger for the module
	Github       github.Client        // A GitHub client for the module
	Blacklist    *blacklist.Blacklist // The blacklist instance for filtering versions
}

// RepositoryURL constructs the URL to the module repository on github.com.
func (m Module) RepositoryURL() string {
	return fmt.Sprintf("https://github.com/%s/terraform-%s-%s", m.Namespace, m.TargetSystem, m.Name)
}

// RSSURL returns the URL of the RSS feed for the repository's tags.
func (m Module) RSSURL() string {
	repositoryURL := m.RepositoryURL()
	return fmt.Sprintf("%s/tags.atom", repositoryURL)
}

// VersionDownloadURL returns the location to download the module from.
// the file should just contain a link to GitHub to download the tarball, ie:
// git::https://github.com/terraform-aws-modules/terraform-aws-iam?ref=<commit-hash>
func (m Module) VersionDownloadURL(version Version) string {
	ref := version.Commit
	if ref == "" {
		// TODO remove this fallback once backfill is complete
		ref = version.Version
	}
	return fmt.Sprintf("git::%s?ref=%s", m.RepositoryURL(), ref)
}

// MetadataPath returns the path to the metadata file for the module.
func (m Module) MetadataPath() string {
	return filepath.Join(m.Directory, strings.ToLower(m.Namespace[0:1]), m.Namespace, m.Name, m.TargetSystem+".json")
}

// ReadMetadata reads the metadata file for the module.
func (m Module) ReadMetadata() (Metadata, error) {
	var metadata Metadata

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

// WriteMetadata writes the metadata to a file.
func (m Module) WriteMetadata(meta Metadata) error {
	path := m.MetadataPath()
	return files.SafeWriteObjectToJSONFile(path, meta)
}
