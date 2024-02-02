package validate

import "golang.org/x/mod/semver"

// IsValidVersion check if the version is semver compatible.
func IsValidVersion(ver string) bool {
	if ver[0] != 'v' {
		ver = "v" + ver
	}
	return semver.IsValid(ver)
}
