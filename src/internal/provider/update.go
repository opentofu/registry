package provider

import (
	"fmt"
	"log/slog"
	"registry-stable/internal/github"

	"golang.org/x/mod/semver"
)

func (p Provider) UpdateMetadataFile() error {
	if shouldUpdate, err := p.shouldUpdateMetadataFile(); err != nil || !shouldUpdate {
		return err
	}

	meta, err := p.buildMetadataFile()
	if err != nil {
		return err
	}
	return p.WriteMetadata(*meta)
}

func (p Provider) shouldUpdateMetadataFile() (bool, error) {
	lastSemverTag, err := p.getLastSemverTag()
	if err != nil {
		return false, err
	}

	fileContent, err := p.ReadMetadata()
	if err != nil {
		return false, err
	}

	for _, v := range fileContent.Versions {
		versionWithPrefix := fmt.Sprintf("v%s", v.Version)
		if versionWithPrefix == lastSemverTag {
			p.Logger.Info("Found latest tag, nothing to update...", slog.String("tag", lastSemverTag))
			return false, nil
		}
	}

	p.Logger.Info("Could not find latest tag, updating...", slog.String("tag", lastSemverTag))
	return true, nil

}

func (p Provider) getSemverTags() ([]string, error) {
	releasesRssUrl := p.getRssUrl()
	tags, err := github.GetTagsFromRss(releasesRssUrl)
	if err != nil {
		return nil, err
	}

	var semverTags []string
	for _, tag := range tags {
		if semver.IsValid(tag) {
			semverTags = append(semverTags, tag)
		}
	}

	return semverTags, nil
}

func (p Provider) getLastSemverTag() (string, error) {
	semverTags, err := p.getSemverTags()
	if err != nil {
		return "", err
	}

	if len(semverTags) < 1 {
		return "", fmt.Errorf("no semver tags found in repository %s", p.RepositoryURL())
	}

	// Tags should be sorted by descending creation date. So, return the first tag
	return semverTags[0], nil
}
