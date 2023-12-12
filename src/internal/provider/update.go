package provider

import (
	"fmt"
	"log/slog"

	"golang.org/x/mod/semver"
)

// UpdateMetadataFile updates the metadata file with the latest version information
func (p Provider) UpdateMetadataFile() error {
	p.Logger.Info("Beginning version bump process")

	meta, err := p.buildMetadata()
	if err != nil {
		p.Logger.Error("Failed to version bump provider", slog.Any("err", err))
		return err
	}
	if meta == nil {
		return nil
	}

	p.Logger.Info("Completed provider version bump successfully")
	return p.WriteMetadata(*meta)
}

func (p Provider) shouldUpdateMetadataFile() (bool, error) {
	semVerTag, err := p.getLastSemVerTag()
	if err != nil {
		return false, err
	}

	if semVerTag == "" {
		// Repo unavailable or tags deleted
		return false, nil
	}

	fileContent, err := p.ReadMetadata()
	if err != nil {
		return false, err
	}

	for _, v := range fileContent.Versions {
		versionWithPrefix := fmt.Sprintf("v%s", v.Version)
		if versionWithPrefix == semVerTag {
			p.Logger.Info("Found latest tag, nothing to update...", slog.String("tag", semVerTag))
			return false, nil
		}
	}

	p.Logger.Info("Could not find latest tag, updating...", slog.String("tag", semVerTag))
	return true, nil
}

// getSemVerTagsFromRSS returns a list of semver tags from the RSS feed
// ignoring all non-valid semver tags
func (p Provider) getSemVerTagsFromRSS() ([]string, error) {
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
func (p Provider) getLastSemVerTag() (string, error) {
	semverTags, err := p.getSemVerTagsFromRSS()
	if err != nil {
		// TODO This is a stopgap, the logs will need to be checked regularly for this.
		p.Logger.Error("Unable to fetch tags, skipping", slog.Any("err", err))
		return "", nil
	}

	if len(semverTags) < 1 {
		// TODO This is a stopgap, the logs will need to be checked regularly for this.
		p.Logger.Error("no semver tags found in repository, skipping", slog.String("url", p.RepositoryURL()))
		return "", nil
	}

	// Tags should be sorted by descending creation date. So, return the first tag
	return semverTags[0], nil
}
