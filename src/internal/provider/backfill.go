package provider

import (
	"context"
	"errors"
	"log/slog"
)

func (p *Provider) BackfillVersionData(ctx context.Context) error {
	p.Logger.Info("Beginning version backfill process")

	meta, err := p.ReadMetadata()
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

		if version.Discovered == nil {
			version.Discovered = new(version.FirstDiscovered())
			meta.Versions[key] = version
			madeChanges = true
			backfilled++
		} else {
			skipped++
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
