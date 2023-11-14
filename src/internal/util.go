package internal

import "strings"

func TrimTagPrefix(version string) string {
	return strings.TrimPrefix(version, "v")
}
