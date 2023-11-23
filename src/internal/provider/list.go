package provider

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"registry-stable/internal/github"
	"registry-stable/internal/parallel"
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
		return nil
	}

	p := Provider{
		Namespace:    matches[providerDirectoryRegex.SubexpIndex("Namespace")],
		ProviderName: matches[providerDirectoryRegex.SubexpIndex("ProviderName")],
	}
	return &p
}

type ProviderList []Provider

func ListProviders(providerDataDir string, logger *slog.Logger, ghClient github.Client) (ProviderList, error) {
	// walk the provider directory recursively and find all json files
	// for each json file, parse it into a Provider struct
	// return a slice of Provider structs

	var results []Provider
	err := filepath.Walk(providerDataDir, func(path string, info os.FileInfo, err error) error {
		p := extractProviderDetailsFromPath(path)
		if p != nil {
			p.Directory = providerDataDir
			p.Logger = logger.With(
				slog.String("type", "provider"),
				slog.Group("provider", slog.String("namespace", p.Namespace), slog.String("name", p.ProviderName)),
			)
			p.Github = ghClient.WithLogger(p.Logger)

			results = append(results, *p)
		} else {
			logger.Debug("Failed to extract provider details from path, skipping", slog.String("path", path))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

type Action func(p Provider) error

func (providers ProviderList) Parallel(maxConcurrency int, action Action) error {
	actions := make([]parallel.Action, len(providers))
	for i, p := range providers {
		p := p
		actions[i] = func() error {
			err := action(p)
			if err != nil {
				p.Logger.Error(err.Error())
				return err
			}
			return nil
		}
	}

	errs := parallel.ForEach(actions, maxConcurrency)
	if len(errs) != 0 {
		return fmt.Errorf("Encountered %d errors processing %d providers", len(errs), len(providers))
	}
	return nil
}
