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

	var errs []error
	madeChanges := false
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
			continue
		}

		newVersion, err := p.VersionFromTag("v" + version.Version)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if newVersion != nil {
			meta.Versions[key] = *newVersion
			madeChanges = true
		}
	}

	p.Logger.Info("Completed version backfill process")

	if madeChanges {
		errs = append(errs, p.WriteMetadata(meta))
	}

	return errors.Join(errs...)
}
