package module

import (
	"fmt"

	"github.com/opentofu/registry-stable/internal/validate"
)

// Validate validates module's metadata.
func Validate(v Metadata) error {
	var errs = make([]error, 0)
	if len(v.Versions) < 1 {
		errs = append(errs, validate.ErrorEmptyList)
	}

	for _, ver := range v.Versions {
		if !validate.IsValidVersion(ver.Version) {
			errs = append(errs, fmt.Errorf("found semver-incompatible version: %s", ver.Version))
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return validate.Errors(errs)
}
