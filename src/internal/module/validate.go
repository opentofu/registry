package module

import (
	"errors"
	"fmt"

	"github.com/opentofu/registry-stable/internal/validate"
)

// Validate validates module's metadata.
func Validate(v Metadata) error {
	if len(v.Versions) < 1 {
		return validate.ErrorEmptyList
	}

	var errs = make([]error, 0, len(v.Versions))
	for _, ver := range v.Versions {
		if !validate.IsValidVersion(ver.Version) {
			err := fmt.Errorf("found semver-incompatible version: %s", ver.Version)
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
