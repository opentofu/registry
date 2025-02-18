package provider

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/opentofu/registry-stable/internal"
)

var (
	// From OpenTOFU go-releaser

	goos = []string{
		"darwin",
		"freebsd",
		"linux",
		"windows",
		"openbsd",
		"solaris",
	}

	goarch = []string{
		"386",
		"amd64",
		"arm",
		"arm64",
	}
)

type VersionFromTagArgs struct {
	Release   string
	URLPrefix string
}

// VersionFromTag fetches information about an individual release based on the GitHub release name
func (p Provider) VersionFromTag(args VersionFromTagArgs) (*Version, error) {
	urlPrefix := args.URLPrefix

	if urlPrefix == "" {
		urlPrefix = p.RepositoryURL() + "/releases/download"
	}

	if args.Release == "" {
		return nil, fmt.Errorf("argument 'release' must be specified")
	}

	release := args.Release
	version := internal.TrimTagPrefix(release)
	lowercaseVersion := strings.ToLower(version)
	artifactPrefix := fmt.Sprintf("%s_%s_", p.RepositoryName(), version)

	logger := p.Logger.With(slog.String("release", release))

	releasePrefix := fmt.Sprintf("%s/%s/%s", urlPrefix, release, artifactPrefix)

	v := Version{
		Version:             lowercaseVersion,
		SHASumsURL:          releasePrefix + "SHA256SUMS",
		SHASumsSignatureURL: releasePrefix + "SHA256SUMS.sig",
	}

	checksums, err := p.GetSHASums(v.SHASumsURL)
	if err != nil {
		return nil, err
	}

	if checksums == nil {
		logger.Warn("checksums not found in release, skipping...")
		return nil, nil
	}

	var ok bool
	for _, os := range goos {
		for _, arch := range goarch {
			target := Target{
				OS:          os,
				Arch:        arch,
				Filename:    fmt.Sprintf("%s%s_%s.zip", artifactPrefix, os, arch),
				DownloadURL: fmt.Sprintf("%s%s_%s.zip", urlPrefix, os, arch),
			}
			target.SHASum, ok = checksums[target.Filename]
			if ok {
				v.Targets = append(v.Targets, target)
				continue
			}
			// now try and pull it with the v in the version
			target.Filename = fmt.Sprintf("%s_v%s_%s_%s.zip", p.RepositoryName(), version, os, arch)
			target.DownloadURL = fmt.Sprintf("%s/v%s/%s", urlPrefix, version, target.Filename)
			target.SHASum, ok = checksums[target.Filename]
			logger.Info("target:", slog.String("target Filename", target.Filename))
			if ok {
				v.Targets = append(v.Targets, target)
			}
		}
	}

	if len(v.Targets) == 0 {
		logger.Info("No artifacts in release, skipping...", slog.String("release", version))
		return nil, nil
	}

	v.Protocols, err = p.GetProtocols(releasePrefix + "manifest.json")
	if err != nil {
		return nil, err
	}
	if v.Protocols == nil {
		logger.Warn("Could not find manifest file, using default protocols")
		v.Protocols = []string{"5.0"}
	}

	return &v, nil
}
