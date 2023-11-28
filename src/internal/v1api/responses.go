package v1api

// ModuleVersionDownloadResponse is the item returned by the module version download API.
type ModuleVersionDownloadResponse struct {
	// The URL to download the module from.
	Location string `json:"location"`
}

// ModuleVersionListingResponse is the item returned by the module version listing API.
type ModuleVersionListingResponse struct {
	Modules []ModuleVersionListingResponseItem `json:"modules"`
}

// ModuleVersionListingResponseItem returns information about all the versions of a module.
type ModuleVersionListingResponseItem struct {
	Versions []ModuleVersionResponseItem `json:"versions"`
}

// ModuleVersionResponseItem is the item returned by the module version listing API
type ModuleVersionResponseItem struct {
	Version string `json:"version"` // The version string

	// TODO: Discuss if we want to keep or not.
	// Root is not currently populated in the response, but may be in the future.
	Root *ModuleMetadata `json:"root,omitempty"`
}

type ModuleMetadata struct {
	Path         string                     `json:"path,omitempty"` // If this is a submodule, the path to the module root
	Providers    []ModuleProviderDependency `json:"providers"`      // The providers required by this module
	Dependencies []ModuleDependency         `json:"dependencies"`   // The modules required by this module

	SubModules []ModuleMetadata `json:"submodules"` // The submodules of this module
}

type ModuleProviderDependency struct {
	Name      string `json:"name"`      // The name of the module
	Namespace string `json:"namespace"` // The namespace of the module
	Version   string `json:"version"`   // The version constraint defined inside the module, ie. ">= 1.0.0"
	Source    string `json:"source"`    // The name of the provider, ie. "hashicorp/aws" or "myregistry.com/myorg/myprovider"
}

type ModuleDependency struct {
	Name    string `json:"name"`    // The name of the module
	Source  string `json:"source"`  // The location of the module
	Version string `json:"version"` // The version constraint defined inside the module, ie. ">= 1.0.0"
}

// ProviderVersionListingResponse is the item returned by the provider version listing API.
type ProviderVersionListingResponse struct {
	Versions []ProviderVersionResponseItem `json:"versions"` // A list of provider versions.
}

type ProviderVersionResponseItem struct {
	Version   string     `json:"version"`   // The version number of the provider.
	Protocols []string   `json:"protocols"` // The protocol versions the provider supports.
	Platforms []Platform `json:"platforms"` // A list of platforms for which this provider version is available.
}

// ProviderVersionDetails provides comprehensive details about a specific provider version.
// This includes the OS, architecture, download URLs, SHA sums, and the signing keys used for the version.
// This is made to match the registry v1 API response format for the download details.
type ProviderVersionDetails struct {
	Protocols           []string    `json:"protocols"`             // The protocol versions the provider supports.
	OS                  string      `json:"os"`                    // The operating system for which the provider is built.
	Arch                string      `json:"arch"`                  // The architecture for which the provider is built.
	Filename            string      `json:"filename"`              // The filename of the provider binary.
	DownloadURL         string      `json:"download_url"`          // The direct URL to download the provider binary.
	SHASumsURL          string      `json:"shasums_url"`           // The URL to the SHA checksums file.
	SHASumsSignatureURL string      `json:"shasums_signature_url"` // The URL to the GPG signature of the SHA checksums file.
	SHASum              string      `json:"shasum"`                // The SHA checksum of the provider binary.
	SigningKeys         SigningKeys `json:"signing_keys"`          // The signing keys used for this provider version.
}

// SigningKeys represents the GPG public keys used to sign a provider version.
type SigningKeys struct {
	GPGPublicKeys []GPGPublicKey `json:"gpg_public_keys"` // A list of GPG public keys.
}

// GPGPublicKey represents an individual GPG public key.
type GPGPublicKey struct {
	KeyID      string `json:"key_id"`      // The ID of the GPG key.
	ASCIIArmor string `json:"ascii_armor"` // The ASCII armored representation of the GPG public key.
}
