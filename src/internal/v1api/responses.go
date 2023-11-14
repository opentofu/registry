package v1api

// ModuleVersionListingResponse is the item returned by the module version listing API.
type ModuleVersionListingResponse struct {
	Modules []ModuleVersionListingResponseItem `json:"modules"`
}

type ModuleVersionListingResponseItem struct {
	Versions []VersionResponseItem `json:"versions"`
}

type VersionResponseItem struct {
	Version string `json:"version"` // The version string

	// Root is not currently populated in the response, but may be in the future.
	Root *ModuleMetadata `json:"root,omitempty"`
}

type ModuleMetadata struct {
	Path         string                     `json:"path,omitempty"` // If this is a submodule, the path to the module root
	Providers    []ModuleProviderDependency `json:"providers"`
	Dependencies []ModuleDependency         `json:"dependencies"`

	SubModules []ModuleMetadata `json:"submodules"`
}

type ModuleProviderDependency struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Version   string `json:"version"` // The version constraint defined inside the module, ie. ">= 1.0.0"
	Source    string `json:"source"`  // The name of the provider, ie. "hashicorp/aws" or "myregistry.com/myorg/myprovider"
}

type ModuleDependency struct {
	Name    string `json:"name"`
	Source  string `json:"source"`
	Version string `json:"version"`
}
