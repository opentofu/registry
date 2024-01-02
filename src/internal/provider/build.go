package provider

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/opentofu/registry-stable/internal"
	"github.com/opentofu/registry-stable/internal/github"
)

// filterNewReleases filters the list of releases to only include those that do
// not already exist in the metadata.
func (p Provider) filterNewReleases(releases []string) []string {
	var existingVersions = make(map[string]bool)
	for _, v := range p.Versions {
		existingVersions[v.Version] = true
	}

	var newReleases = make([]string, 0)
	for _, r := range releases {
		if !existingVersions[internal.TrimTagPrefix(r)] {
			// only append the release if it does not already exist in the metadata
			newReleases = append(newReleases, r)
		}
	}

	p.Log.Info(fmt.Sprintf("Found %d releases that do not already exist in the metadata file", len(newReleases)))

	return newReleases
}

func (p *Provider) UpdateMetadata() error {
	// fetch ALL the releases
	tags, err := p.Repository.ListTags()
	if err != nil {
		p.Log.Error("Unable to fetch semver tags, skipping", slog.Any("err", err))
		return nil
	}

	// Filter to semver onlu
	tags = tags.FilterSemver()

	// filter the releases to only include those that do not already exist in the metadata
	newReleases := p.filterNewReleases(tags)

	if len(newReleases) == 0 {
		p.Log.Info("No version bump required, all versions exist")
		return nil
	}

	shouldUpdate, err := p.shouldUpdateMetadataFile()
	if err != nil {
		p.Log.Error("Failed to determine update status, forcing update", slog.Any("err", err))
	} else if !shouldUpdate {
		p.Log.Info("No version bump required, latest versions exist")
		return nil
	}

	type versionResult struct {
		v   *Version
		err error
	}

	verChan := make(chan versionResult, len(newReleases))

	// for each of the new releases, fetch the version and add it to the metadata
	for _, r := range newReleases {
		r := r
		go func() {
			version, err := p.VersionFromTag(r)
			verChan <- versionResult{version, err}
		}()
	}

	// TODO agg errors
	for range newReleases {
		result := <-verChan
		if result.err != nil {
			return result.err
		}
		if result.v == nil {
			// Not a valid release, skipping
			continue
		}
		// append the new release to the metadata
		p.Versions = append(p.Versions, *result.v)
	}

	slices.SortFunc(p.Versions, func(a, b Version) int {
		return github.SemverTagSort(a.Version, b.Version)
	})

	return nil
}
