package module

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/opentofu/registry-stable/internal"
	"github.com/opentofu/registry-stable/internal/github"

	"golang.org/x/mod/semver"
)

func (m Module) UpdateMetadataFile() error {
	m.Logger.Info("Beginning version bump process for module", slog.String("module", m.Namespace+"/"+m.Name+"/"+m.TargetSystem))

	meta, err := m.BuildMetadata()
	if err != nil {
		return err
	}
	if meta == nil {
		return nil
	}

	return m.WriteMetadata(*meta)
}

// BuildMetadata builds the Metadata for the module by collating the tags from
// the module repository.
func (m Module) BuildMetadata() (*Metadata, error) {
	tags, err := m.getSemverTags()
	if err != nil {
		m.Logger.Error("Unable to fetch semver tags, skipping", slog.Any("err", err))
		return nil, nil
	}

	meta, err := m.ReadMetadata()
	if err != nil {
		return nil, err
	}

	// Merge current versions with new versions
	for _, t := range tags {
		found := false
		for _, v := range meta.Versions {
			if v.Version == t.Name {
				found = true
				break
			}
		}
		if !found {
			meta.Versions = append(meta.Versions, Version{Version: t.Name, Ref: t.Commit})
		}
	}

	semverSortFunc := func(a, b Version) int {
		return -semver.Compare(fmt.Sprintf("v%s", internal.TrimTagPrefix(a.Version)), fmt.Sprintf("v%s", internal.TrimTagPrefix(b.Version)))
	}
	slices.SortFunc(meta.Versions, semverSortFunc)

	return &meta, nil
}

func (m Module) getSemverTags() ([]github.Tag, error) {
	tags, err := m.Github.GetTags(m.RepositoryURL())
	if err != nil {
		return nil, err
	}

	var semverTags = make([]github.Tag, 0)
	for _, tag := range tags {
		tagWithPrefix := fmt.Sprintf("v%s", internal.TrimTagPrefix(tag.Name))
		if semver.IsValid(tagWithPrefix) {
			semverTags = append(semverTags, tag)
		}
	}

	semverSortFunc := func(a, b github.Tag) int {
		return -semver.Compare(a.Name, b.Name)
	}
	slices.SortFunc(semverTags, semverSortFunc)

	return semverTags, nil
}
