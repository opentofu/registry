package provider

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/opentofu/registry-stable/internal/blacklist"
	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/github"
)

// Metadata contains information about the provider.
type Metadata struct {
	Repository string       `json:"repository,omitempty"` // Optional. Custom repository from which to fetch the provider's metadata.
	Versions   []Version    `json:"versions"`             // A list of version data, for each supported provider version.
	Warnings   []string     `json:"warnings,omitempty"`   // Warnings associated with this provider.
	Logger     *slog.Logger `json:"-"`
}

// Version contains information about a specific provider version.
type Version struct {
	Version             string   `json:"version"`               // The version number of the provider.
	Protocols           []string `json:"protocols"`             // The protocol versions the provider supports.
	SHASumsURL          string   `json:"shasums_url"`           // The URL to the SHA checksums file.
	SHASumsSignatureURL string   `json:"shasums_signature_url"` // The URL to the GPG signature of the SHA checksums file.
	Targets             []Target `json:"targets"`               // A list of target platforms for which this provider version is available.
}

// Target contains information about a specific provider version for a specific target platform.
type Target struct {
	OS          string `json:"os"`           // The operating system for which the provider is built.
	Arch        string `json:"arch"`         // The architecture for which the provider is built.
	Filename    string `json:"filename"`     // The filename of the provider release.
	DownloadURL string `json:"download_url"` // The direct URL to download the provider release.
	SHASum      string `json:"shasum"`       // The SHA checksum of the provider release.
	Hash1       string `json:"h1"`           // The Hash Type 1 of the provider release *contents*
}

// Provider contains information about a provider.
type Provider struct {
	ProviderName string // The provider name
	Namespace    string // The provider namespace
	Directory    string // The root directory that the provider lives in
	Logger       *slog.Logger
	Github       github.Client
	Blacklist    *blacklist.Blacklist // The blacklist instance for filtering versions
}

// RepositoryName returns the name of the repository that the provider is assumed to be living in.
func (p Provider) RepositoryName() string {
	return fmt.Sprintf("terraform-provider-%s", p.ProviderName)
}

// RepositoryURL returns the URL of the repository that the provider is assumed to be living in.
func (p Provider) RepositoryURL() string {
	return fmt.Sprintf("https://github.com/%s/%s", p.EffectiveNamespace(), p.RepositoryName())
}

// RSSURL returns the URL of the RSS feed for the repository's releases.
func (p Provider) RSSURL() string {
	repositoryUrl := p.RepositoryURL()
	return fmt.Sprintf("%s/releases.atom", repositoryUrl)
}

// EffectiveNamespace will map namespaces for providers in situations
// where the author (owner of the namespace) does not release artifacts as
// GitHub Releases.
func (p Provider) EffectiveNamespace() string {
	if p.Namespace == "hashicorp" {
		return "opentofu"
	}

	return p.Namespace
} // TODO make more generic

// MetadataPath returns the path to the provider's metadata file.
func (p Provider) MetadataPath() string {
	return filepath.Join(p.Directory, strings.ToLower(p.Namespace[0:1]), p.Namespace, p.ProviderName+".json")
}

// ReadMetadata reads and parses the provider's metadata file and returns a Metadata struct.
func (p Provider) ReadMetadata() (Metadata, error) {
	var metadata Metadata

	path := p.MetadataPath()

	metadataFile, err := os.ReadFile(path)
	if err != nil {
		return metadata, fmt.Errorf("failed to open metadata file: %w", err)
	}

	err = json.Unmarshal(metadataFile, &metadata)
	if err != nil {
		return metadata, fmt.Errorf("failed to unmarshal metadata file: %w", err)
	}

	metadata.Logger = p.Logger

	return metadata, nil
}

// WriteMetadata writes the given Metadata struct to the provider's metadata file.
func (p Provider) WriteMetadata(meta Metadata) error {
	path := p.MetadataPath()
	return files.SafeWriteObjectToJSONFile(path, meta)
}
