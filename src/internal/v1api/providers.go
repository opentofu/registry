package v1api

import (
	"context"
	"fmt"
	"path/filepath"

	"registry-stable/internal/files"
	"registry-stable/internal/provider"
)

// GenerateProviderResponses generates the response for the provider version listing API endpoints.
func (g Generator) GenerateProviderResponses(_ context.Context, p provider.Provider) error {
	metadata, err := p.ReadMetadata()
	if err != nil {
		return err
	}

	s := ProviderSource{p, metadata}

	for location, details := range s.VersionDetails() {
		path := filepath.Join(g.DestinationDir, location)
		err := files.SafeWriteObjectToJsonFile(path, details)
		if err != nil {
			return fmt.Errorf("failed to write metadata version download file: %w", err)
		}
	}

	path := filepath.Join(g.DestinationDir, s.VersionListingPath())
	err = files.SafeWriteObjectToJsonFile(path, s.Versions())
	if err != nil {
		return err
	}

	return nil
}
