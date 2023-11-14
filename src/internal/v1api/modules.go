package v1api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"registry-stable/internal"
	"registry-stable/internal/module"
)

// GenerateModuleResponses generates the response for the module version listing API endpoints.
// For more information see
// https://opentofu.org/docs/internals/module-registry-protocol/#list-available-versions-for-a-specific-module
// https://opentofu.org/docs/internals/module-registry-protocol/#download-source-code-for-a-specific-module-version
func (g Generator) GenerateModuleResponses(_ context.Context, namespace string, name string, targetSystem string) error {
	logger := slog.With(slog.String("namespace", namespace), slog.String("name", name), slog.String("targetSystem", targetSystem))

	// TODO: Get path calculation from somewhere else
	path := filepath.Join(g.ModuleDirectory, namespace[0:1], namespace, name, targetSystem+".json")

	metadata, err := g.readModuleMetadata(path, logger)
	if err != nil {
		return err
	}

	// Right now the format is pretty much identical, however if we want to extend the results in the future to include
	// more information, we can do that here. (i.e. the root module, or the submodules)
	versionsResponse := make([]ModuleVersionResponseItem, len(metadata.Versions))
	for i, m := range metadata.Versions {
		versionsResponse[i] = ModuleVersionResponseItem{Version: m.Version}

		err := g.writeModuleVersionDownload(namespace, name, targetSystem, m.Version)
		if err != nil {
			return fmt.Errorf("failed to write metadata version download file for version %s: %w", m.Version, err)
		}
		logger.Debug("Wrote metadata version download file", slog.String("version", m.Version))
	}

	// Write the /versions response
	err = g.writeModuleVersionListing(namespace, name, targetSystem, versionsResponse)
	if err != nil {
		return err
	}

	return nil
}

// readModuleMetadata reads the module metadata file from the filesystem directly. This data should be the data fetched from the git repository.
func (g Generator) readModuleMetadata(path string, logger *slog.Logger) (*module.MetadataFile, error) {
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

	logger.Debug("Loaded Modules", slog.Any("count", len(metadata.Versions)))

	return &metadata, nil
}

// writeModuleVersionListing writes the file containing the module version listing.
// This data  is to be consumed when an end user requests /v1/modules/{namespace}/{name}/{targetSystem}/versions
func (g Generator) writeModuleVersionListing(namespace string, name string, targetSystem string, versions []ModuleVersionResponseItem) error {
	destinationDir := filepath.Join(g.DestinationDir, "v1", "modules", namespace, name, targetSystem)
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(destinationDir, "versions")

	marshalled, err := json.Marshal(ModuleVersionListingResponse{Modules: []ModuleVersionListingResponseItem{{Versions: versions}}})
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	err = os.WriteFile(filePath, marshalled, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// writeModuleVersionDownload writes the file containing the download link for the module version.
// This data is to be consumed when an end user requests /v1/modules/{namespace}/{name}/{targetSystem}/{version}/download
func (g Generator) writeModuleVersionDownload(namespace string, name string, system string, version string) interface{} {
	// the file should just contain a link to GitHub to download the tarball, ie:
	// git::https://github.com/terraform-aws-modules/terraform-aws-iam?ref=v5.30.0
	location := fmt.Sprintf("git::github.com/%s/terraform-%s-%s?ref=%s", namespace, name, system, version)

	response := ModuleVersionDownloadResponse{Location: location}

	marshalled, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	// trim the v from the version as the api only ever requests the version without the v
	destinationDir := filepath.Join(g.DestinationDir, "v1", "modules", namespace, name, system, internal.TrimTagPrefix(version))
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(destinationDir, "download")
	err = os.WriteFile(filePath, marshalled, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
