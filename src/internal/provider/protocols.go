package provider

import (
	"encoding/json"
	"fmt"
)

type Manifest struct {
	Metadata ManifestMetadata `json:"metadata"`
}
type ManifestMetadata struct {
	ProtocolVersions []string `json:"protocol_versions"`
}

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
		return nil, err
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
