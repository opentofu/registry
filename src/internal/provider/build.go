package provider

import (
	"fmt"
	"slices"

	"registry-stable/internal"
	"registry-stable/internal/github"

	"golang.org/x/mod/semver"
)

func (meta MetadataFile) filterNewReleases(releases []github.GHRelease) []github.GHRelease {
	var existingVersions = make(map[string]bool)
	for _, v := range meta.Versions {
		existingVersions[v.Version] = true
	}

	var newReleases = make([]github.GHRelease, 0)
	for _, r := range releases {
		if !existingVersions[internal.TrimTagPrefix(r.TagName)] {
			newReleases = append(newReleases, r)
		}
	}

	meta.Logger.Info(fmt.Sprintf("Found %d releases that do not already exist in the metadata file", len(newReleases)))

	return newReleases
}

func (p Provider) buildMetadataFile() (*MetadataFile, error) {
	meta, err := p.ReadMetadata()
	if err != nil {
		return nil, err
	}

	releases, err := p.Github.FetchPublishedReleases(p.EffectiveNamespace(), p.RepositoryName())
	if err != nil {
		return nil, err
	}

	releases = meta.filterNewReleases(releases)

	for _, r := range releases {
		version, err := p.VersionFromRelease(r)
		if err != nil {
			return nil, err
		}
		if version == nil {
			// Not a valid release, skipping
			continue
		}

		meta.Versions = append(meta.Versions, *version)
	}

	semverSortFunc := func(a, b Version) int {
		return semver.Compare(fmt.Sprintf("s%s", a.Version), fmt.Sprintf("s%s", b.Version))
	}
	slices.SortFunc(meta.Versions, semverSortFunc)

	return &meta, nil
}
