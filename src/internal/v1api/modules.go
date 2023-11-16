package v1api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"registry-stable/internal/files"
	"registry-stable/internal/module"
)

// GenerateModuleResponses generates the response for the module version listing API endpoints.
// For more information see
// https://opentofu.org/docs/internals/module-registry-protocol/#list-available-versions-for-a-specific-module
// https://opentofu.org/docs/internals/module-registry-protocol/#download-source-code-for-a-specific-module-version
func (g Generator) GenerateModuleResponses(_ context.Context, m module.Module) error {
	logger := slog.With(slog.String("namespace", m.Namespace), slog.String("name", m.Name), slog.String("targetSystem", m.TargetSystem))

	metadata, err := m.ReadMetadata(g.ModuleDirectory, logger)
	if err != nil {
		return err
	}

	// Right now the format is pretty much identical, however if we want to extend the results in the future to include
	// more information, we can do that here. (i.e. the root module, or the submodules)
	versionsResponse := make([]ModuleVersionResponseItem, len(metadata.Versions))
	for i, v := range metadata.Versions {
		versionsResponse[i] = ModuleVersionResponseItem{Version: v.Version}

		err := g.writeModuleVersionDownload(m, v)
		if err != nil {
			return fmt.Errorf("failed to write metadata version download file for version %s: %w", v.Version, err)
		}
		logger.Debug("Wrote metadata version download file", slog.String("version", v.Version))
	}

	// Write the /versions response
	err = g.writeModuleVersionListing(m, versionsResponse)
	if err != nil {
		return err
	}

	return nil
}

// writeModuleVersionListing writes the file containing the module version listing.
// This data  is to be consumed when an end user requests /v1/modules/{namespace}/{name}/{targetSystem}/versions
func (g Generator) writeModuleVersionListing(m module.Module, versions []ModuleVersionResponseItem) error {
	return files.SafeWriteObjectToJsonFile(
		filepath.Join(g.DestinationDir, m.VersionListingPath()),
		ModuleVersionListingResponse{Modules: []ModuleVersionListingResponseItem{{Versions: versions}}},
	)
}

// readModuleMetadata reads the module metadata file from the filesystem directly. This data should be the data fetched from the git repository.
func (g Generator) readModuleMetadata(m module.Module, logger *slog.Logger) (*module.MetadataFile, error) {
	path := filepath.Join(g.ModuleDirectory, m.MetadataPath())

	// open the file
	metadataFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}

	// Read the file contents into a Module[] struct
	var metadata module.MetadataFile
	err = json.Unmarshal(metadataFile, &metadata)
	if err != nil {
		return nil, err
	}

	logger.Debug("Loaded Module Versions", slog.Any("count", len(metadata.Versions)))

	return &metadata, nil
}

// writeModuleVersionDownload writes the file containing the download link for the module version.
// This data is to be consumed when an end user requests /v1/modules/{namespace}/{name}/{targetSystem}/{version}/download
func (g Generator) writeModuleVersionDownload(m module.Module, version module.Version) interface{} {
	return files.SafeWriteObjectToJsonFile(
		filepath.Join(g.DestinationDir, m.VersionDownloadPath(version)),
		ModuleVersionDownloadResponse{Location: m.VersionDownloadURL(version)},
	)
}
