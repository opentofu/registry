package provider

import (
	"context"
	"fmt"
	"io"
	"registry-stable/internal/github"
	"strings"
)

func GetShaSums(ctx context.Context, shaFileDownloadUrl string) (map[string]string, error) {
	assetContents, assetErr := github.DownloadAssetContents(ctx, shaFileDownloadUrl)
	if assetErr != nil {
		return nil, fmt.Errorf("failed to download asset contents: %w", assetErr)
	}

	contents, contentsErr := io.ReadAll(assetContents)
	if contentsErr != nil {
		return nil, fmt.Errorf("failed to read asset contents: %w", contentsErr)
	}

	return shaFileToMap(contents), nil
}

func shaFileToMap(contents []byte) map[string]string {
	var result = make(map[string]string)
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		shaSum := parts[0]
		fileName := parts[1]

		result[fileName] = shaSum
	}
	return result
}
