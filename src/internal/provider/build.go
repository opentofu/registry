package provider

import (
	"context"
	"fmt"
	"log"
	"slices"

	"registry-stable/internal"
	"registry-stable/internal/github"

	"golang.org/x/mod/semver"
)

func (existingMetadata MetadataFile) filterNewReleases(releases []github.GHRelease) []github.GHRelease {
	var existingVersions = make(map[string]bool)
	for _, v := range existingMetadata.Versions {
		existingVersions[v.Version] = true
	}

	var newReleases = make([]github.GHRelease, 0)
	for _, r := range releases {
		if !existingVersions[internal.TrimTagPrefix(r.TagName)] {
			newReleases = append(newReleases, r)
		}
	}

	log.Printf("Found %d releases that do not already exist in the metadata file", len(newReleases))

	return newReleases
}

func (p Provider) buildMetadataFile() (*MetadataFile, error) {
	ctx := context.Background()

	token, err := github.EnvAuthToken()
	if err != nil {
		return nil, err
	}

	ghClient := github.NewGitHubClient(token)

	meta, err := p.ReadMetadata()
	if err != nil {
		return nil, err
	}

	releases, err := github.FetchPublishedReleases(ctx, ghClient, p.EffectiveNamespace(), p.RepositoryName())
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
