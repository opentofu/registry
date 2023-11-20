package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"registry-stable/internal/files"
	"strings"
)

type MetadataFile struct {
	Repository string    `json:"repository,omitempty"` // Optional. Custom repository from which to fetch the provider's metadata.
	Versions   []Version `json:"versions"`             // A list of version data, for each supported provider version.
}

type Version struct {
	Version             string   `json:"version"`               // The version number of the provider.
	Protocols           []string `json:"protocols"`             // The protocol versions the provider supports.
	SHASumsURL          string   `json:"shasums_url"`           // The URL to the SHA checksums file.
	SHASumsSignatureURL string   `json:"shasums_signature_url"` // The URL to the GPG signature of the SHA checksums file.
	Targets             []Target `json:"targets"`               // A list of target platforms for which this provider version is available.
}

type Target struct {
	OS          string `json:"os"`           // The operating system for which the provider is built.
	Arch        string `json:"arch"`         // The architecture for which the provider is built.
	Filename    string `json:"filename"`     // The filename of the provider binary.
	DownloadURL string `json:"download_url"` // The direct URL to download the provider binary.
	SHASum      string `json:"shasum"`       // The SHA checksum of the provider binary.
}

type Provider struct {
	ProviderName string // The provider name
	Namespace    string // The provider namespace
	Directory    string // The root directory that the provider lives in
}

// TODO remove me and use slog instead?
func (p Provider) String() string {
	return fmt.Sprintf("%s/%s", p.ProviderName, p.Namespace)
}

func (p Provider) RepositoryName() string {
	return fmt.Sprintf("terraform-provider-%s", p.ProviderName)
}

// TODO custom repository url?
func (p Provider) RepositoryURL() string {
	return fmt.Sprintf("https://github.com/%s/%s", p.EffectiveNamespace(), p.RepositoryName())
}

func (p Provider) getRssUrl() string {
	repositoryUrl := p.RepositoryURL()
	return fmt.Sprintf("%s/releases.atom", repositoryUrl)
}

// EffectiveProviderNamespace will map namespaces for providers in situations
// where the author (owner of the namespace) does not release artifacts as
// GitHub Releases.
func (p Provider) EffectiveNamespace() string {
	if p.Namespace == "hashicorp" {
		return "opentofu"
	}

	return p.Namespace
} // TODO make more generic

func (p Provider) MetadataPath() string {
	return filepath.Join(p.Directory, strings.ToLower(p.Namespace[0:1]), p.Namespace, p.ProviderName+".json")
}

func (p Provider) ReadMetadata() (MetadataFile, error) {
	var metadata MetadataFile

	path := p.MetadataPath()

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

func (p Provider) WriteMetadata(meta MetadataFile) error {
	path := p.MetadataPath()
	return files.SafeWriteObjectToJsonFile(path, meta)
}
