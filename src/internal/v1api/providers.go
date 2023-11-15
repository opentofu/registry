package v1api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"registry-stable/internal/files"
	"registry-stable/internal/github"
	"registry-stable/internal/provider"
)

// GenerateProviderResponses generates the response for the provider version listing API endpoints.
func (g Generator) GenerateProviderResponses(_ context.Context, namespace string, name string) error {
	logger := slog.With(slog.String("namespace", namespace), slog.String("name", name))

	path := filepath.Join(g.ProviderDirectory, strings.ToLower(namespace[0:1]), namespace, name+".json")

	metadata, err := g.readProviderMetadata(path, logger)
	if err != nil {
		return err
	}

	versionsResponse := make([]ProviderVersionResponseItem, len(metadata.Versions))
	for versionIdx, ver := range metadata.Versions {
		// TODO: Extract to a nice copy constructor method, rather than doing it inline here
		versionsResponse[versionIdx] = ProviderVersionResponseItem{Version: ver.Version}

		// construct the Platforms from the `Targets` in the metadata
		platforms := make([]github.Platform, len(ver.Targets))

		versionDetails := []ProviderVersionDetails{}

		for targetIdx, target := range ver.Targets {
			platforms[targetIdx] = github.Platform{
				OS:   target.OS,
				Arch: target.Arch,
			}

			versionDetails = append(versionDetails, ProviderVersionDetails{
				Protocols:           ver.Protocols,
				OS:                  target.OS,
				Arch:                target.Arch,
				Filename:            target.Filename,
				DownloadURL:         target.DownloadURL,
				SHASumsURL:          ver.SHASumsURL,
				SHASumsSignatureURL: ver.SHASumsSignatureURL,
				SHASum:              target.SHASum,
				SigningKeys:         SigningKeys{}, // TODO: Add gpg keys
			})
		}
		versionsResponse[versionIdx].Platforms = platforms
		versionsResponse[versionIdx].Protocols = ver.Protocols

		// for each of the targets, write the version download file
		for _, details := range versionDetails {
			logger.Debug("Writing version download file", slog.String("version", ver.Version))
			err := g.writeProviderVersionDownload(namespace, name, ver.Version, details)
			if err != nil {
				return fmt.Errorf("failed to write metadata version download file for version %s: %w", ver.Version, err)
			}
		}
		logger.Debug("Wrote metadata version download file", slog.String("version", ver.Version))
	}

	err = g.writeProviderVersionListing(namespace, name, versionsResponse)
	if err != nil {
		return err
	}

	return nil
}

// readProviderMetadata reads the provider metadata file from the filesystem directly. This data should be the data fetched from the git repository.
func (g Generator) readProviderMetadata(path string, logger *slog.Logger) (*provider.MetadataFile, error) {
	metadataFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}

	var metadata provider.MetadataFile
	err = json.Unmarshal(metadataFile, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata file: %w", err)
	}

	logger.Debug("Loaded Provider Metadata", slog.Any("versions", len(metadata.Versions)))

	return &metadata, nil
}

func (g Generator) writeProviderVersionDownload(namespace string, name string, version string, versionMetadata ProviderVersionDetails) error {
	path := filepath.Join(g.DestinationDir, "v1", "providers", namespace, name, version, "download", versionMetadata.OS, versionMetadata.Arch)
	return files.WriteToJsonFile(path, versionMetadata)
}

func (g Generator) writeProviderVersionListing(namespace string, name string, versions []ProviderVersionResponseItem) error {
	path := filepath.Join(g.DestinationDir, "v1", "providers", namespace, name, "versions")
	return files.WriteToJsonFile(path, ProviderVersionListingResponse{Versions: versions})
}
