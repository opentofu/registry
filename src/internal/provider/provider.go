package provider

import (
	"fmt"
	"log/slog"

	"github.com/opentofu/registry-stable/internal/github"
)

// Metadata contains information about the provider.
type Metadata struct {
	Repository string    `json:"repository,omitempty"` // Optional. Custom repository from which to fetch the provider's metadata.
	Versions   []Version `json:"versions"`             // A list of version data, for each supported provider version.
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
	Filename    string `json:"filename"`     // The filename of the provider binary.
	DownloadURL string `json:"download_url"` // The direct URL to download the provider binary.
	SHASum      string `json:"shasum"`       // The SHA checksum of the provider binary.
}

type Identifier struct {
	ProviderName string // The provider name
	Namespace    string // The provider namespace
}

func (id Identifier) String() string {
	return fmt.Sprintf("%s/%s", id.Namespace, id.ProviderName)
}

// Provider contains information about a provider.
type Provider struct {
	Identifier
	Metadata
	Log        *slog.Logger
	Repository github.Repository
}
