package provider

import (
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
)

/*
providerDirectoryRegex is a regular expression that matches the directory structure of a provider file.
  - (?i) makes the match case-insensitive.
  - providers/ matches the literal string "providers/".
  - \w matches a single word character. This corresponds to the first letter of the namespace.
  - (?P<Namespace>[^/]+) captures a sequence of one or more characters that are not a slash. This corresponds to "oracle".
  - (?P<ProviderName>[^/]+) captures another sequence of non-slash characters. This corresponds to "oci".
  - \.json matches the literal string ".json".
*/
var providerDirectoryRegex = regexp.MustCompile(`(?i)providers/\w/(?P<Namespace>[^/]+?)/(?P<ProviderName>[^/]+?)\.json`)

func extractProviderDetailsFromPath(path string) *Provider {
	matches := providerDirectoryRegex.FindStringSubmatch(path)
	if len(matches) != 3 {
		slog.Debug("Failed to extract provider details from path, skipping", slog.String("path", path))
		return nil
	}

	p := Provider{
		Namespace:    matches[providerDirectoryRegex.SubexpIndex("Namespace")],
		ProviderName: matches[providerDirectoryRegex.SubexpIndex("ProviderName")],
	}

	return &p
}

func ListProviders(providerDataDir string) ([]Provider, error) {
	// walk the provider directory recursively and find all json files
	// for each json file, parse it into a Provider struct
	// return a slice of Provider structs

	var results []Provider
	err := filepath.Walk(providerDataDir, func(path string, info os.FileInfo, err error) error {
		p := extractProviderDetailsFromPath(path)
		if p != nil {
			results = append(results, *p)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
