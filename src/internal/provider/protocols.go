package provider

import (
	"encoding/json"
	"fmt"
)

// Manifest contains information about the provider manifest.
type Manifest struct {
	Metadata ManifestMetadata `json:"metadata"`
}

// ManifestMetadata contains information about the provider manifest metadata.
type ManifestMetadata struct {
	ProtocolVersions []string `json:"protocol_versions"`
}

func parseManifestContents(contents []byte) (*Manifest, error) {
	var manifest *Manifest
	err := json.Unmarshal(contents, &manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifest contents: %w", err)
	}

	return manifest, nil
}
