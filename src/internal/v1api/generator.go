package v1api

// Generator is responsible for generating the responses used by v1 registry APIs.
// This should take information from the current state of the registry and generate the responses
// to be served by the API.
// For more information on the API see the following documentation:
// - https://opentofu.org/docs/internals/module-registry-protocol/
// - https://opentofu.org/docs/internals/provider-registry-protocol
type Generator struct {
	// DestinationDir is the directory to write the generated responses to.
	DestinationDir string

	// ModuleDirectory is the directory to read module metadata from.
	ModuleDirectory string

	// ProviderDirectory is the directory to read provider metadata from.
	ProviderDirectory string
}
