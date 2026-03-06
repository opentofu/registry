package provider

import (
	"context"
	"errors"
)

func (p *Provider) BackfillVersionData(ctx context.Context) error {
	p.Logger.Info("Beginning version backfill process")

	meta, err := p.ReadMetadata()
	if err != nil {
		return err
	}

	madeChanges := false
	for key, version := range meta.Versions {
		if err := ctx.Err(); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				// Outta-time
				break
			}
			return err
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
			continue
		}

		newVersion, err := p.VersionFromTag("v" + version.Version)
		if err != nil {
			// TODO log
			continue
		}
		if newVersion != nil {
			meta.Versions[key] = *newVersion
			madeChanges = true
		}
	}

	p.Logger.Info("Completed version backfill process")

	if madeChanges {
		return p.WriteMetadata(meta)
	}

	return nil
}
