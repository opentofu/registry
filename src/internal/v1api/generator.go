package v1api

import (
	"io/fs"

	"registry-stable/internal/files"
)

// Generator is responsible for generating the responses used by v1 registry APIs.
// This should take information from the current state of the registry and generate the responses
// to be served by the API.
// For more information on the API see the following documentation:
// - https://opentofu.org/docs/internals/module-registry-protocol/
// - https://opentofu.org/docs/internals/provider-registry-protocol
type Generator struct {
	// DestinationDir is the directory to write the generated responses to.
	DestinationDir string

	// TODO: Maybe combine all the fs.FS stuff and FileWriter stuff into a single interface to be consumed easier?

	// ModuleFS is the filesystem the generator should use to read files for moduleMetadata
	ModuleFS fs.FS

	// FileWriter is to be used to write files to the filesystem
	FileWriter files.FileWriter
}
