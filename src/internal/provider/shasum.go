package provider

import (
	"fmt"
	"strings"
)

// GetSHASums will attempt to download the SHA checksums file from the given URL and return a
// map of file names to SHA checksums.
func (p Provider) GetSHASums(shaFileDownloadUrl string) (map[string]string, error) {
	contents, assetErr := p.Github.DownloadAssetContents(shaFileDownloadUrl)
	if assetErr != nil {
		return nil, fmt.Errorf("failed to download asset contents: %w", assetErr)
	}
	if contents == nil {
		return nil, nil
	}

	return shaFileToMap(contents), nil
}

// shaFileToMap converts the contents of a SHA checksums file into a map of file names to SHA checksums.
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
