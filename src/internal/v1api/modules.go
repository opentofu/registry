package v1api

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/opentofu/registry-stable/internal"
	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/module"
)

// ModuleGenerator is responsible for generating the response for the module version listing API endpoints.
type ModuleGenerator struct {
	module.Module
	module.Metadata
	Destination string
	log         *slog.Logger
}

// NewModuleGenerator creates a new ModuleGenerator which will generate the response for the module version listing API endpoints
// and store the generated files in the given destination directory
func NewModuleGenerator(m module.Module, destination string) (ModuleGenerator, error) {
	metadata, err := m.ReadMetadata()
	if err != nil {
		return ModuleGenerator{}, err
	}

	return ModuleGenerator{
		Module:      m,
		Metadata:    metadata,
		Destination: destination,
		log:         m.Logger,
	}, nil
}

// VersionListingPath returns the path to the module version listing file
func (m ModuleGenerator) VersionListingPath() string {
	namespace := strings.ToLower(m.Module.Namespace)
	name := strings.ToLower(m.Module.Name)
	target := strings.ToLower(m.Module.TargetSystem)
	return filepath.Join(m.Destination, "v1", "modules", namespace, name, target, "versions")

}

// VersionDownloadPath returns the path to the module version download file for the given version
func (m ModuleGenerator) VersionDownloadPath(v module.Version) string {
	namespace := strings.ToLower(m.Module.Namespace)
	name := strings.ToLower(m.Module.Name)
	target := strings.ToLower(m.Module.TargetSystem)
	version := strings.ToLower(internal.TrimTagPrefix(v.Version))
	return filepath.Join(m.Destination, "v1", "modules", namespace, name, target, version, "download")
}

// VersionListing converts the module metadata into a ModuleVersionListingResponse, ready to be serialized to a file
func (m ModuleGenerator) VersionListing() ModuleVersionListingResponse {
	versions := make([]ModuleVersionResponseItem, len(m.Metadata.Versions))
	for i, v := range m.Metadata.Versions {
		versions[i] = ModuleVersionResponseItem{Version: internal.TrimTagPrefix(v.Version)}
	}
	return ModuleVersionListingResponse{[]ModuleVersionListingResponseItem{{versions}}}
}

// VersionDownloads converts the module metadata into a map of module version download paths to ModuleVersionDownloadResponse,
func (m ModuleGenerator) VersionDownloads() map[string]ModuleVersionDownloadResponse {
	downloads := make(map[string]ModuleVersionDownloadResponse)
	for _, v := range m.Metadata.Versions {
		downloads[m.VersionDownloadPath(v)] = ModuleVersionDownloadResponse{Location: m.Module.VersionDownloadURL(v)}
	}
	return downloads
}

// Generate generates the response for the module version listing API endpoints.
// For more information see
// https://opentofu.org/docs/internals/module-registry-protocol/#list-available-versions-for-a-specific-module
// https://opentofu.org/docs/internals/module-registry-protocol/#download-source-code-for-a-specific-module-version
func (m ModuleGenerator) Generate() error {
	m.log.Info("Generating")

	for location, download := range m.VersionDownloads() {
		err := files.SafeWriteObjectToJSONFile(location, download)
		if err != nil {
			return fmt.Errorf("failed to write metadata version download file: %w", err)
		}
		m.log.Debug("Wrote metadata version download file", slog.String("path", location))
	}

	err := files.SafeWriteObjectToJSONFile(m.VersionListingPath(), m.VersionListing())
	if err != nil {
		return err
	}

	m.log.Info("Generated")

	return nil
}
