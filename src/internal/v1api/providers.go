package v1api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"registry-stable/internal/files"
	"registry-stable/internal/github"
	"registry-stable/internal/provider"
)

// GenerateProviderResponses generates the response for the provider version listing API endpoints.
func (g Generator) GenerateProviderResponses(_ context.Context, p provider.Provider) error {
	logger := slog.With(slog.String("namespace", p.Namespace), slog.String("name", p.ProviderName))

	metadata, err := g.readProviderMetadata(p, logger)
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
			err := g.writeProviderVersionDownload(p, ver, details)
			if err != nil {
				return fmt.Errorf("failed to write metadata version download file for version %s: %w", ver.Version, err)
			}
		}
		logger.Debug("Wrote metadata version download file", slog.String("version", ver.Version))
	}

	err = g.writeProviderVersionListing(p, versionsResponse)
	if err != nil {
		return err
	}

	return nil
}

// readProviderMetadata reads the provider metadata file from the filesystem directly. This data should be the data fetched from the git repository.
func (g Generator) readProviderMetadata(p provider.Provider, logger *slog.Logger) (*provider.MetadataFile, error) {
	path := filepath.Join(g.ProviderDirectory, p.MetadataPath())

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

func (g Generator) writeProviderVersionDownload(p provider.Provider, v provider.Version, versionMetadata ProviderVersionDetails) error {
	path := filepath.Join(g.DestinationDir, "v1", p.VersionDownloadPath(v), versionMetadata.OS, versionMetadata.Arch)
	return files.SafeWriteObjectToJsonFile(path, versionMetadata)
}

func (g Generator) writeProviderVersionListing(p provider.Provider, versions []ProviderVersionResponseItem) error {
	path := filepath.Join(g.DestinationDir, "v1", p.VersionListingPath())
	return files.SafeWriteObjectToJsonFile(path, ProviderVersionListingResponse{Versions: versions})
}
