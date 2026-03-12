package provider

import (
	"context"
	"errors"
	"log/slog"
	"slices"
)

func (p *Provider) BackfillVersionData(ctx context.Context) error {
	p.Logger.Info("Beginning version backfill process")

	meta, err := p.ReadMetadata()
	if err != nil {
		return err
	}

	releases, err := p.getSemverTags()
	if err != nil {
		return err
	}

	var errs []error
	madeChanges := false
	var backfilled, skipped, errored int
	for key, version := range meta.Versions {
		if err := ctx.Err(); err != nil {
			errs = append(errs, err)
			break
		}

		// Check to see if this version has already
		// been populated
		hasMeta := false
		for _, target := range version.Targets {
			if target.Size != 0 {
				hasMeta = true
				break
			}
		}
		if hasMeta {
			skipped++
			continue
		}

		// version.Version *never* starts with "v"
		// Therefore we need to find the git tag which corresponds
		// with the recorded provider version.
		var releaseTag string
		if slices.Contains(releases, version.Version) {
			releaseTag = version.Version
		} else if slices.Contains(releases, "v"+version.Version) {
			releaseTag = "v" + version.Version
		} else {
			p.Logger.Warn("Failed to backfill version, missing release tag", slog.String("version", version.Version))
			break
		}

		newVersion, err := p.VersionFromTag(releaseTag)
		if err != nil {
			p.Logger.Error("Failed to backfill version", slog.String("version", version.Version), slog.Any("err", err))
			errs = append(errs, err)
			errored++
			continue
		}
		if newVersion != nil {
			meta.Versions[key] = *newVersion
			madeChanges = true
			backfilled++
		}
	}

	p.Logger.Info("Completed version backfill process",
		slog.Int("backfilled", backfilled),
		slog.Int("skipped", skipped),
		slog.Int("errored", errored),
		slog.Int("total", len(meta.Versions)),
	)

	if madeChanges {
		errs = append(errs, p.WriteMetadata(meta))
	}

	return errors.Join(errs...)
}
