package github

import (
	"regexp"
)

// TODO: probably move the Platform type inside providers as that's the only place this is used,
// and its a bit strange being in the github package

type Platform struct {
	OS   string
	Arch string
}

func ExtractPlatformFromFilename(filename string) *Platform {
	platformPattern := regexp.MustCompile(`.*_(?P<Os>[a-zA-Z0-9]+)_(?P<Arch>[a-zA-Z0-9]+).zip`)
	matches := platformPattern.FindStringSubmatch(filename)

	if matches == nil {
		return nil
	}

	platform := Platform{
		OS:   matches[platformPattern.SubexpIndex("Os")],
		Arch: matches[platformPattern.SubexpIndex("Arch")],
	}

	return &platform
}
