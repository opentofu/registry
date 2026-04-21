package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

func (m *Module) BackfillVersionData(ctx context.Context) error {
	m.Logger.Info("Beginning version backfill process")

	meta, err := m.ReadMetadata()
	if err != nil {
		return err
	}

	releases, err := m.getSemverTags()
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

		if version.Discovered != nil && version.Commit != "" {
			skipped++
			continue
		}
		madeChanges = true

		if version.Discovered == nil {
			version.Discovered = new(version.FirstDiscovered())
			meta.Versions[key] = version
		}

		if version.Commit == "" {
			foundRelease := false
			for _, release := range releases {
				if release.Ref == version.Version {
					version.Commit = release.Commit
					meta.Versions[key] = version
					foundRelease = true
					break
				}
			}
			if !foundRelease {
				err := fmt.Errorf("release not found for version %s", version.Version)
				m.Logger.Error("Failed to backfill version", slog.String("version", version.Version), slog.Any("err", err))
				errs = append(errs, err)
				errored++
			}
		}

		backfilled++
	}

	m.Logger.Info("Completed version backfill process",
		slog.Int("backfilled", backfilled),
		slog.Int("skipped", skipped),
		slog.Int("errored", errored),
		slog.Int("total", len(meta.Versions)),
	)

	if madeChanges {
		errs = append(errs, m.WriteMetadata(meta))
	}

	return errors.Join(errs...)
}
