package v1api

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/opentofu/registry-stable/internal"
	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/module"
)

// ModuleGenerator is responsible for generating the response for the module version listing API endpoints.
type ModuleGenerator struct {
	module.Module
	Destination string
}

// NewModuleGenerator creates a new ModuleGenerator which will generate the response for the module version listing API endpoints
// and store the generated files in the given destination directory
func NewModuleGenerator(m module.Module, destination string) ModuleGenerator {
	return ModuleGenerator{
		Module:      m,
		Destination: destination,
	}
}

// VersionListingPath returns the path to the module version listing file
func (m ModuleGenerator) VersionListingPath() string {
	return filepath.Join(m.Destination, "v1", "modules", m.Namespace, m.Name, m.TargetSystem, "versions")
}

// VersionDownloadPath returns the path to the module version download file for the given version
func (m ModuleGenerator) VersionDownloadPath(v module.Version) string {
	return filepath.Join(m.Destination, "v1", "modules", m.Namespace, m.Name, m.TargetSystem, internal.TrimTagPrefix(v.Version), "download")
}

// VersionListing converts the module metadata into a ModuleVersionListingResponse, ready to be serialized to a file
func (m ModuleGenerator) VersionListing() ModuleVersionListingResponse {
	versions := make([]ModuleVersionResponseItem, len(m.Versions))
	for i, v := range m.Versions {
		versions[i] = ModuleVersionResponseItem{Version: internal.TrimTagPrefix(v.Version)}
	}
	return ModuleVersionListingResponse{[]ModuleVersionListingResponseItem{{versions}}}
}

// VersionDownloads converts the module metadata into a map of module version download paths to ModuleVersionDownloadResponse,
func (m ModuleGenerator) VersionDownloads() map[string]ModuleVersionDownloadResponse {
	downloads := make(map[string]ModuleVersionDownloadResponse)
	for _, v := range m.Versions {
		downloads[m.VersionDownloadPath(v)] = ModuleVersionDownloadResponse{Location: m.Repository.DownloadURL(v.Version)}
	}
	return downloads
}

// Generate generates the response for the module version listing API endpoints.
// For more information see
// https://opentofu.org/docs/internals/module-registry-protocol/#list-available-versions-for-a-specific-module
// https://opentofu.org/docs/internals/module-registry-protocol/#download-source-code-for-a-specific-module-version
func (m ModuleGenerator) Generate() error {
	m.Log.Info("Generating")

	for location, download := range m.VersionDownloads() {
		err := files.SafeWriteObjectToJSONFile(location, download)
		if err != nil {
			return fmt.Errorf("failed to write metadata version download file: %w", err)
		}
		m.Log.Debug("Wrote metadata version download file", slog.String("path", location))
	}

	err := files.SafeWriteObjectToJSONFile(m.VersionListingPath(), m.VersionListing())
	if err != nil {
		return err
	}

	m.Log.Info("Generated")

	return nil
}
