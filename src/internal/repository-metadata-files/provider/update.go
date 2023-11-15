package provider

import (
	"encoding/json"
	"fmt"
	"golang.org/x/mod/semver"
	"log"
	"os"
	"registry-stable/internal/github"
	"registry-stable/internal/provider"
)

func UpdateMetadataFile(p provider.Provider) error {
	if shouldUpdate, err := shouldUpdateMetadataFile(p); err != nil || !shouldUpdate {
		return err
	}

	return CreateMetadataFile(p)
}

func shouldUpdateMetadataFile(p provider.Provider) (bool, error) {
	lastSemverTag, err := getLastSemverTag(p)
	if err != nil {
		return false, err
	}

	pathToFile := getFilePath(p)
	fileContent, err := getProviderFileContent(pathToFile)
	if err != nil {
		return false, err
	}

	for _, v := range fileContent.Versions {
		versionWithPrefix := fmt.Sprintf("v%s", v.Version)
		if versionWithPrefix == lastSemverTag {
			log.Printf("Found latest tag %s in the repository file %s, nothing to update...", lastSemverTag, pathToFile)
			return false, nil
		}
	}

	log.Printf("Could not find latest tag %s in the repository file %s, updating the file...", lastSemverTag, pathToFile)
	return true, nil

}

func getProviderFileContent(path string) (provider.MetadataFile, error) {
	res, _ := os.ReadFile(path)

	var fileData provider.MetadataFile

	err := json.Unmarshal(res, &fileData)

	if err != nil {
		return provider.MetadataFile{}, err
	}

	return fileData, nil
}

func getSemverTags(p provider.Provider) ([]string, error) {
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

func getLastSemverTag(p provider.Provider) (string, error) {
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

func getRssUrl(p provider.Provider) string {
	repositoryUrl := p.RepositoryURL()
	return fmt.Sprintf("%s/releases.atom", repositoryUrl)
}
