package github

import (
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

type Tags []string

func (tags Tags) FilterSemver() Tags {
	filtered := make(Tags, 0, len(tags))
	for _, tag := range tags {
		if semver.IsValid(SemverFormat(tag)) {
			filtered = append(filtered, tag)
		}
	}
	return filtered
}

func SemverFormat(tag string) string {
	return fmt.Sprintf("v%s", strings.TrimPrefix(tag, "v"))
}

func SemverTagSort(a, b string) int {
	// TODO remove the inversion
	return -semver.Compare(SemverFormat(a), SemverFormat(b))
}
