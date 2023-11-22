package module

import (
	"fmt"
	"registry-stable/internal"

	"golang.org/x/mod/semver"
)

func (m Module) UpdateMetadataFile() error {
	meta, err := m.BuildMetadataFile()
	if err != nil {
		return err
	}

	return m.WriteMetadata(*meta)
}

func (m Module) BuildMetadataFile() (*MetadataFile, error) {
	tags, err := m.getSemverTags()
	if err != nil {
		return nil, err
	}

	var versions = make([]Version, 0, len(tags))
	for _, t := range tags {
		versions = append(versions, Version{Version: t})
	}

	return &MetadataFile{Versions: versions}, nil
}

func (m Module) getSemverTags() ([]string, error) {
	tags, err := m.Github.GetTags(m.RepositoryURL())
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
