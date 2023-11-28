package module

import (
	"fmt"
	"log/slog"

	"golang.org/x/mod/semver"
)

func (p Module) shouldUpdateMetadataFile() (bool, error) {
	semVerTag, err := p.getLastSemVerTag()
	if err != nil {
		return false, err
	}

	fileContent, err := p.ReadMetadata()
	if err != nil {
		return false, err
	}

	for _, v := range fileContent.Versions {
		if v.Version == semVerTag {
			p.Logger.Info("Found latest tag, nothing to update...", slog.String("tag", semVerTag))
			return false, nil
		}
	}

	p.Logger.Info("Could not find latest tag, updating...", slog.String("tag", semVerTag))
	return true, nil
}

// getSemVerTagsFromRSS returns a list of semver tags from the RSS feed
// ignoring all non-valid semver tags
func (p Module) getSemVerTagsFromRSS() ([]string, error) {
	releasesRssUrl := p.RSSURL()
	tags, err := p.Github.GetTagsFromRSS(releasesRssUrl)
	if err != nil {
		return nil, err
	}

	var semverTags []string
	for _, tag := range tags {
		if semver.IsValid(tag) || semver.IsValid("v"+tag) {
			semverTags = append(semverTags, tag)
		}
	}

	return semverTags, nil
}

// getLastSemVerTag returns the most recently created semver tag from the RSS feed
// by sorting the tags by descending creation date
func (p Module) getLastSemVerTag() (string, error) {
	semverTags, err := p.getSemVerTagsFromRSS()
	if err != nil {
		return "", err
	}

	if len(semverTags) < 1 {
		return "", fmt.Errorf("no semver tags found in repository %s", p.RepositoryURL())
	}

	// Tags should be sorted by descending creation date. So, return the first tag
	return semverTags[0], nil
}
