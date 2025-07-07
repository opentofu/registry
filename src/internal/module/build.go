package module

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/opentofu/registry-stable/internal"
	"github.com/opentofu/registry-stable/internal/blacklist"

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
	// Use blacklist from Module struct, fall back to empty if nil
	blacklistInstance := m.Blacklist
	if blacklistInstance == nil {
		blacklistInstance = &blacklist.Blacklist{}
	}

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
	blacklistedCount := 0
	for _, t := range tags {
		found := false
		for _, v := range meta.Versions {
			if v.Version == t {
				found = true
				break
			}
		}
		if !found {
			// Check if this version is blacklisted
			version := internal.TrimTagPrefix(t)
			if isBlacklisted, reason := blacklistInstance.IsModuleVersionBlacklisted(m.Namespace, m.Name, m.TargetSystem, version); isBlacklisted {
				m.Logger.Warn("Skipping blacklisted module version", 
					slog.String("namespace", m.Namespace),
					slog.String("name", m.Name),
					slog.String("target", m.TargetSystem),
					slog.String("version", version),
					slog.String("reason", reason))
				blacklistedCount++
				continue
			}
			meta.Versions = append(meta.Versions, Version{Version: t})
		}
	}
	
	if blacklistedCount > 0 {
		m.Logger.Info(fmt.Sprintf("Skipped %d blacklisted versions", blacklistedCount))
	}

	semverSortFunc := func(a, b Version) int {
		return -semver.Compare(fmt.Sprintf("v%s", internal.TrimTagPrefix(a.Version)), fmt.Sprintf("v%s", internal.TrimTagPrefix(b.Version)))
	}
	slices.SortFunc(meta.Versions, semverSortFunc)

	return &meta, nil
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

	semverSortFunc := func(a, b string) int {
		return -semver.Compare(fmt.Sprintf(a), fmt.Sprintf(b))
	}
	slices.SortFunc(semverTags, semverSortFunc)

	return semverTags, nil
}
