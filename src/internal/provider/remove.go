package provider

import "fmt"

// RemoveVersion removes the given version from the provider's metadata and
// writes the updated metadata back to disk
func (p Provider) RemoveVersion(version string) error {
	meta, err := p.ReadMetadata()
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	filtered := make([]Version, 0, len(meta.Versions))
	found := false
	for _, v := range meta.Versions {
		if v.Version == version {
			found = true
			continue
		}
		filtered = append(filtered, v)
	}

	if !found {
		return fmt.Errorf("version %q not found in metadata for %s/%s", version, p.Namespace, p.ProviderName)
	}

	meta.Versions = filtered

	if err := p.WriteMetadata(meta); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}
