package internal

import "strings"

// TODO: Move this somewhere to it's own package, not in internal.

func TrimTagPrefix(version string) string {
	return strings.TrimPrefix(version, "v")
}
