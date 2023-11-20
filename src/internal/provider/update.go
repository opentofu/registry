package provider

import (
	"fmt"
	"log"
	"registry-stable/internal/github"

	"golang.org/x/mod/semver"
)

func UpdateMetadataFile(p Provider, providerDataDir string) error {
	if shouldUpdate, err := shouldUpdateMetadataFile(p, providerDataDir); err != nil || !shouldUpdate {
		return err
	}

	meta, err := BuildMetadataFile(p, providerDataDir)
	if err != nil {
		return err
	}
	return p.WriteMetadata(providerDataDir, *meta)
}

func shouldUpdateMetadataFile(p Provider, providerDataDir string) (bool, error) {
	lastSemverTag, err := getLastSemverTag(p)
	if err != nil {
		return false, err
	}

	fileContent, err := p.ReadMetadata(providerDataDir)
	if err != nil {
		return false, err
	}

	for _, v := range fileContent.Versions {
		versionWithPrefix := fmt.Sprintf("v%s", v.Version)
		if versionWithPrefix == lastSemverTag {
			log.Printf("Found latest tag %s for %s, nothing to update...", lastSemverTag, p.String())
			return false, nil
		}
	}

	log.Printf("Could not find latest tag %s for %s, updating the file...", lastSemverTag, p.String())
	return true, nil

}

func getSemverTags(p Provider) ([]string, error) {
	releasesRssUrl := getRssUrl(p)
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

func getLastSemverTag(p Provider) (string, error) {
	semverTags, err := getSemverTags(p)
	if err != nil {
		return "", err
	}

	if len(semverTags) < 1 {
		return "", fmt.Errorf("no semver tags found in repository for provider %s/%s", p.Namespace, p.ProviderName)
	}

	// Tags should be sorted by descending creation date. So, return the first tag
	return semverTags[0], nil
}

func getRssUrl(p Provider) string {
	repositoryUrl := p.RepositoryURL()
	return fmt.Sprintf("%s/releases.atom", repositoryUrl)
}
