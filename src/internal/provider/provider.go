package provider

import "fmt"

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
}

func (p Provider) RepositoryName() string {
	return fmt.Sprintf("terraform-provider-%s", p.ProviderName)
}

// TODO custom repository url?
func (p Provider) RepositoryURL() string {
	return fmt.Sprintf("https://github.com/%s/%s", p.EffectiveNamespace(), p.RepositoryName())
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
