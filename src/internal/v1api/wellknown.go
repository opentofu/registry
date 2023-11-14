package v1api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const WellKnownFileContents = `{
	  "modules.v1": "/v1/modules/",
	  "providers.v1": "/v1/providers/"
}`

// WriteWellKnownFile writes the well-known file to the filesystem.
// For more information see
// https://opentofu.org/docs/internals/remote-service-discovery/#discovery-process
func (g Generator) WriteWellKnownFile(_ context.Context) error {
	wellKnownDir := filepath.Join(g.DestinationDir, ".well-known")
	err := os.MkdirAll(wellKnownDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(wellKnownDir, "terraform.json")
	err = os.WriteFile(filePath, []byte(WellKnownFileContents), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
