package provider

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/opentofu/registry-stable/internal"

	"golang.org/x/mod/semver"
)

// filterNewReleases filters the list of releases to only include those that do
// not already exist in the metadata.
func (meta Metadata) filterNewReleases(releases []string) []string {
	var existingVersions = make(map[string]bool)
	for _, v := range meta.Versions {
		existingVersions[v.Version] = true
	}

	var newReleases = make([]string, 0)
	for _, r := range releases {
		if !existingVersions[internal.TrimTagPrefix(r)] {
			// only append the release if it does not already exist in the metadata
			newReleases = append(newReleases, r)
		}
	}

	meta.Logger.Info(fmt.Sprintf("Found %d releases that do not already exist in the metadata file", len(newReleases)))

	return newReleases
}

// getSemverTags returns a list of semver tags for the module fetched from GitHub.
func (p Provider) getSemverTags() ([]string, error) {
	tags, err := p.Github.GetTags(p.RepositoryURL())
	if err != nil {
		return nil, err
	}

	var semverTags = make([]string, 0)
	for _, tag := range tags {
		tagWithPrefix := fmt.Sprintf("v%s", internal.TrimTagPrefix(tag.Name))
		if semver.IsValid(tagWithPrefix) {
			semverTags = append(semverTags, tag.Name)
		}
	}

	return semverTags, nil
}

func (p Provider) buildMetadata() (*Metadata, error) {
	meta, err := p.ReadMetadata()
	if err != nil {
		return nil, err
	}

	// fetch ALL the releases
	releases, err := p.getSemverTags()
	if err != nil {
		p.Logger.Error("Unable to fetch semver tags, skipping", slog.Any("err", err))
		return nil, nil
	}

	// filter the releases to only include those that do not already exist in the metadata
	newReleases := meta.filterNewReleases(releases)

	if len(newReleases) == 0 {
		p.Logger.Info("No version bump required, all versions exist")
		return nil, nil
	}

	shouldUpdate, err := p.shouldUpdateMetadataFile()
	if err != nil {
		p.Logger.Error("Failed to determine update status, forcing update", slog.Any("err", err))
	} else if !shouldUpdate {
		p.Logger.Info("No version bump required, latest versions exist")
		return nil, nil
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

	for range newReleases {
		result := <-verChan
		if result.err != nil {
			return nil, result.err
		}
		if result.v == nil {
			// Not a valid release, skipping
			continue
		}
		// append the new release to the metadata
		meta.Versions = append(meta.Versions, *result.v)
	}

	semverSortFunc := func(a, b Version) int {
		return -semver.Compare(fmt.Sprintf("v%s", a.Version), fmt.Sprintf("v%s", b.Version))
	}
	slices.SortFunc(meta.Versions, semverSortFunc)

	return &meta, nil
}
