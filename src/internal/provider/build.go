package provider

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"slices"
	"time"

	"github.com/opentofu/registry-stable/internal"
	"github.com/opentofu/registry-stable/internal/blacklist"

	"golang.org/x/mod/semver"
)

// filterNewReleases filters the list of releases to only include those that do
// not already exist in the metadata and are not blacklisted.
func (meta Metadata) filterNewReleases(releases []string, namespace, name string, blacklistInstance *blacklist.Blacklist) []string {
	var existingVersions = make(map[string]bool)
	for _, v := range meta.Versions {
		existingVersions[v.Version] = true
	}

	var newReleases = make([]string, 0)
	var blacklistedCount = 0
	var erroredCount = 0
	for _, r := range releases {
		version := internal.TrimTagPrefix(r)
		if !existingVersions[version] {
			// Check if this version is blacklisted
			if isBlacklisted, reason := blacklistInstance.IsProviderVersionBlacklisted(namespace, name, version); isBlacklisted {
				meta.Logger.Warn("Skipping blacklisted version",
					slog.String("namespace", namespace),
					slog.String("name", name),
					slog.String("version", version),
					slog.String("reason", reason))
				blacklistedCount++
				continue
			}

			var numVersionErrors int
			var latestError time.Time
			for _, errored := range meta.VersionErrors {
				if errored.Version == version {
					if errored.UTCTime.After(latestError) {
						latestError = errored.UTCTime
					}
					numVersionErrors += 1
				}
			}

			isErrored := false
			if numVersionErrors > 0 {
				// Simple doubling backoff
				dur := time.Minute * 15 * time.Duration(math.Pow(float64(numVersionErrors), 2))
				if latestError.Add(dur).After(time.Now()) {
					isErrored = true
				}
			}

			if isErrored {
				meta.Logger.Warn("Skipping errored version",
					slog.String("namespace", namespace),
					slog.String("name", name),
					slog.String("version", version))
				erroredCount++
				continue
			}

			// only append the release if it does not already exist in the metadata, is not errored, and is not blacklisted
			newReleases = append(newReleases, r)
		}
	}

	meta.Logger.Info(fmt.Sprintf("Found %d releases that do not already exist in the metadata file (%d blacklisted, %d errored)", len(newReleases), blacklistedCount, erroredCount))

	return newReleases
}

// getSemverTags returns a list of semver tags for the module fetched from GitHub.
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

func (p Provider) buildMetadata() (*Metadata, error) {
	meta, err := p.ReadMetadata()
	if err != nil {
		return nil, err
	}
	blacklistInstance := p.Blacklist

	// fetch ALL the releases
	releases, err := p.getSemverTags()
	if err != nil {
		p.Logger.Warn("Unable to fetch semver tags, skipping", slog.Any("err", err))
		return nil, nil
	}

	// filter the releases to only include those that do not already exist in the metadata
	newReleases := meta.filterNewReleases(releases, p.Namespace, p.ProviderName, blacklistInstance)

	if len(newReleases) == 0 {
		p.Logger.Info("No version bump required, all versions exist")
		return nil, nil
	}

	type versionResult struct {
		r   string
		v   *Version
		err error
	}

	verChan := make(chan versionResult, len(newReleases))

	// for each of the new releases, fetch the version and add it to the metadata
	for _, r := range newReleases {
		r := r
		go func() {
			version, err := p.VersionFromTag(r)
			verChan <- versionResult{r, version, err}
		}()
	}

	metadataUpdated := false
	var errs []error
	for range newReleases {
		result := <-verChan
		if result.err != nil {
			var nonFatal ErrVersionNonFatal
			if errors.As(result.err, &nonFatal) {
				meta.VersionErrors = append(meta.VersionErrors, VersionError{
					Version: internal.TrimTagPrefix(result.r),
					Message: nonFatal.Error(),
					UTCTime: time.Now().UTC(),
				})
				metadataUpdated = true
				continue
			}
			errs = append(errs, result.err)
			continue
		}
		if result.v == nil {
			// Not a valid release, skipping
			continue
		}
		// append the new release to the metadata
		meta.Versions = append(meta.Versions, *result.v)

		// remove any error version records now that we have had a successful result
		meta.VersionErrors = slices.DeleteFunc(meta.VersionErrors, func(errored VersionError) bool {
			return errored.Version == result.r
		})

		metadataUpdated = true
	}
	if !metadataUpdated {
		// Prevent file modification if not changed
		return nil, errors.Join(errs...)
	}

	semverSortFunc := func(a, b Version) int {
		return -semver.Compare(fmt.Sprintf("v%s", a.Version), fmt.Sprintf("v%s", b.Version))
	}
	slices.SortFunc(meta.Versions, semverSortFunc)

	slices.SortFunc(meta.VersionErrors, func(a, b VersionError) int {
		comp := -semver.Compare(fmt.Sprintf("v%s", a.Version), fmt.Sprintf("v%s", b.Version))
		if comp == 0 {
			comp = a.UTCTime.Compare(b.UTCTime)
		}
		return comp
	})

	return &meta, errors.Join(errs...)
}
