package provider

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/opentofu/registry-stable/internal/validate"
)

// Validate validates provider's metadata.
func Validate(v Metadata) error {
	if len(v.Versions) < 1 {
		return validate.ErrorEmptyList
	}

	var errs = make([]error, 0, len(v.Versions))
	for _, ver := range v.Versions {
		for _, err := range validateProviderVersion(ver) {
			errs = append(errs, fmt.Errorf("v%s: %w", ver.Version, err))
		}
	}

	return errors.Join(errs...)
}

func validateProviderVersion(ver Version) []error {
	var errs = make([]error, 0)

	if !validate.IsValidVersion(ver.Version) {
		errs = append(errs, fmt.Errorf("found semver-incompatible version: %s", ver.Version))
	}

	// validates provider's protocols:
	if len(ver.Protocols) == 0 {
		errs = append(errs, fmt.Errorf("empty protocols list"))
	} else {
		for _, protocol := range ver.Protocols {
			if !isValidProviderProtocol(protocol) {
				errs = append(errs, fmt.Errorf("unsupported protocol found: %s", protocol))
			}
		}
	}

	// validates provider's targets:
	if len(ver.Targets) == 0 {
		errs = append(errs, fmt.Errorf("empty targets list"))
	} else {
		for _, verTarget := range ver.Targets {
			if err := validateProviderVersionTarget(verTarget); err != nil {
				errs = append(errs, err...)
			}
		}
	}

	return errs
}

// isValidProviderProtocol validates the protocol version.
// It's based on the providers which are currently available in the registry.
func isValidProviderProtocol(s string) bool {
	switch s {
	case "1.0",
		"1.0.0",
		"4.0",
		"5.0",
		"6.0":
		return true
	default:
		return false
	}
}

func validateProviderVersionTarget(v Target) []error {
	var errs = make([]error, 0)

	if !slices.Contains(goos, v.OS) {
		errs = append(errs, fmt.Errorf("target %s-%s: unsupported OS: %s", v.OS, v.Arch, v.OS))
	}

	if !slices.Contains(goarch, v.Arch) {
		errs = append(errs, fmt.Errorf("target %s-%s: unsupported ARCH: %s", v.OS, v.Arch, v.Arch))
	}

	// check if the filename matches the url
	if !strings.HasSuffix(v.DownloadURL, v.Filename) {
		errs = append(errs, fmt.Errorf("target %s-%s: 'filename' is not consistent with 'download_url'", v.OS, v.Arch))
	}

	// check if the SHA sum was modified
	if len(v.SHASum) != 64 {
		errs = append(errs, fmt.Errorf("target %s-%s: SHASum length is wrong", v.OS, v.Arch))
	}

	return errs
}
