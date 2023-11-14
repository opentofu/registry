package module

import (
	"fmt"
	"registry-stable/internal"
	"registry-stable/internal/github"
	"registry-stable/internal/module"

	"golang.org/x/mod/semver"
)

func BuildMetadataFile(m module.Module) (*module.MetadataFile, error) {
	tags, err := getModuleSemverTags(m)
	if err != nil {
		return nil, err
	}

	var versions = make([]module.Version, 0)
	for _, t := range tags {
		versions = append(versions, module.Version{Version: t})
	}

	return &module.MetadataFile{Versions: versions}, nil
}

func getModuleSemverTags(mod module.Module) ([]string, error) {
	tags, err := github.GetTags(mod.RepositoryURL())
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
