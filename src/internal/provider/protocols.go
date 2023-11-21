package provider

import (
	"context"
	"encoding/json"
	"log/slog"
	"registry-stable/internal/github"
)

type Manifest struct {
	Metadata ManifestMetadata `json:"metadata"`
}
type ManifestMetadata struct {
	ProtocolVersions []string `json:"protocol_versions"`
}

func (p Provider) GetProtocols(ctx context.Context, manifestDownloadUrl string) ([]string, error) {
	contents, err := github.DownloadAssetContents(ctx, p.Logger, manifestDownloadUrl)

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
		slog.Error("Failed to parse manifest contents")
		return nil, err
	}

	return manifest, nil
}
