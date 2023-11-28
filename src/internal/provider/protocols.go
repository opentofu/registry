package provider

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

// Manifest contains information about the provider manifest.
type Manifest struct {
	Metadata ManifestMetadata `json:"metadata"`
}

// ManifestMetadata contains information about the provider manifest metadata.
type ManifestMetadata struct {
	ProtocolVersions []string `json:"protocol_versions"`
}

// GetProtocols will attempt to download the manifest from the given URL and return the
// list of protocols that the provider supports.
func (p Provider) GetProtocols(manifestDownloadUrl string) ([]string, error) {
	contents, err := p.Github.DownloadAssetContents(manifestDownloadUrl)
	if err != nil {
		return nil, err
	}
	if contents == nil {
		return nil, nil
	}

	manifest, err := parseManifestContents(contents)
	if err != nil {
		p.Logger.Warn("Manifest file invalid, ignoring...", slog.String("url", manifestDownloadUrl), slog.Any("err", err))
		return nil, nil
	}

	return manifest.Metadata.ProtocolVersions, nil
}

func parseManifestContents(contents []byte) (*Manifest, error) {
	var manifest *Manifest
	err := json.Unmarshal(contents, &manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifest contents: %w", err)
	}

	return manifest, nil
}
