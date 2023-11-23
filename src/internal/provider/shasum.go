package provider

import (
	"fmt"
	"strings"
)

func (p Provider) GetShaSums(shaFileDownloadUrl string) (map[string]string, error) {
	contents, assetErr := p.Github.DownloadAssetContents(shaFileDownloadUrl)
	if assetErr != nil {
		return nil, fmt.Errorf("failed to download asset contents: %w", assetErr)
	}
	if contents == nil {
		return nil, nil
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
