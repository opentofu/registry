package provider

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/opentofu/registry-stable/internal/blacklist"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/parallel"
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

// ListProviders returns a slice of Provider structs for each provider found in the providerDataDir.
func ListProviders(providerDataDir string, providerNamespace string, logger *slog.Logger, ghClient github.Client, blacklistInstance *blacklist.Blacklist) (List, error) {
	var results []Provider
	providerNamespace = strings.ToLower(providerNamespace)
	err := filepath.Walk(providerDataDir, func(path string, info os.FileInfo, err error) error {
		p := extractProviderDetailsFromPath(path)

		if p == nil {
			logger.Debug("Failed to extract provider details from path, skipping", slog.String("path", path))
			return nil
		}

		if providerNamespace != "" && !strings.HasPrefix(strings.ToLower(p.Namespace), providerNamespace) {
			logger.Debug("Filtered out provider", slog.String("path", path))
			return nil
		}

		// enrich the provider object
		p.Directory = providerDataDir
		p.Logger = logger.With(
			slog.String("type", "provider"),
			slog.Group("provider", slog.String("namespace", p.Namespace), slog.String("name", p.ProviderName)),
		)
		p.Github = ghClient.WithLogger(p.Logger)
		p.Blacklist = blacklistInstance

		results = append(results, *p)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// List is a slice of Provider structs.
type List []Provider

// Action is a function that takes a Provider and returns an error, to be used with List.Parallel.
type Action func(p Provider) error

func (providers List) Parallel(maxConcurrency int, action Action) error {
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
		return fmt.Errorf("encountered %d errors processing %d providers", len(errs), len(providers))
	}
	return nil
}
