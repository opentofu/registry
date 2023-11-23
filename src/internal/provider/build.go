package provider

import (
	"fmt"
	"slices"

	"registry-stable/internal"

	"golang.org/x/mod/semver"
)

func (meta MetadataFile) filterNewReleases(releases []string) []string {
	var existingVersions = make(map[string]bool)
	for _, v := range meta.Versions {
		existingVersions[v.Version] = true
	}

	var newReleases = make([]string, 0)
	for _, r := range releases {
		if !existingVersions[internal.TrimTagPrefix(r)] {
			newReleases = append(newReleases, r)
		}
	}

	meta.Logger.Info(fmt.Sprintf("Found %d releases that do not already exist in the metadata file", len(newReleases)))

	return newReleases
}

func (p Provider) getSemverTags() ([]string, error) {
	tags, err := p.Github.GetTags(p.RepositoryURL())
	if err != nil {
		return nil, err
	}

	var semverTags = make([]string, 0)
	for _, tag := range tags {
		tagWithPrefix := fmt.Sprintf("v%s", internal.TrimTagPrefix(tag))
		if semver.IsValid(tagWithPrefix) {
			semverTags = append(semverTags, tag)
		}
	}

	return semverTags, nil
}

func (p Provider) buildMetadataFile() (*MetadataFile, error) {
	meta, err := p.ReadMetadata()
	if err != nil {
		return nil, err
	}

	releases, err := p.getSemverTags()
	if err != nil {
		return nil, err
	}

	releases = meta.filterNewReleases(releases)

	type versionResult struct {
		v   *Version
		err error
	}

	verChan := make(chan versionResult, len(releases))

	for _, r := range releases {
		r := r
		go func() {
			version, err := p.VersionFromTag(r)
			verChan <- versionResult{version, err}
		}()
	}

	for _ = range releases {
		result := <-verChan
		if result.err != nil {
			return nil, result.err
		}
		if result.v == nil {
			// Not a valid release, skipping
			continue
		}
		meta.Versions = append(meta.Versions, *result.v)
	}

	semverSortFunc := func(a, b Version) int {
		return -semver.Compare(fmt.Sprintf("v%s", a.Version), fmt.Sprintf("v%s", b.Version))
	}
	slices.SortFunc(meta.Versions, semverSortFunc)

	return &meta, nil
}
