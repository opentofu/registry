package provider

import (
	"errors"
)

// UpdateMetadataFile updates the metadata file with the latest version information
func (p Provider) UpdateMetadataFile() error {
	p.Logger.Info("Beginning version bump process")

	meta, err := p.buildMetadata()
	if meta != nil {
		err = errors.Join(err, p.WriteMetadata(*meta))
	}

	p.Logger.Info("Completed provider version bump")

	return err
}
