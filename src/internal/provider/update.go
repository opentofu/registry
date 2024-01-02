package provider

import (
	"log/slog"

	"github.com/opentofu/registry-stable/internal/github"
)

func (p Provider) shouldUpdateMetadataFile() (bool, error) {
	tags, err := p.Repository.GetLatestReleases()
	if err != nil {
		return false, err
	}

	tags = tags.FilterSemver()

	if len(tags) == 0 {
		p.Log.Warn("no semver tags found in repository, skipping")
		return false, nil
	}
	tag := tags[0]

	for _, v := range p.Metadata.Versions {
		if github.SemverFormat(v.Version) == github.SemverFormat(tag) {
			p.Log.Info("Found latest tag, nothing to update...", slog.String("tag", tag))
			return false, nil
		}
	}

	p.Log.Info("Could not find latest tag, updating...", slog.String("tag", tag))
	return true, nil
}
