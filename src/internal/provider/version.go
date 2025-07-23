package provider

import (
	"fmt"
	"log/slog"

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

// VersionFromTag fetches information about an individual release based on the GitHub release name
func (p Provider) VersionFromTag(release string) (*Version, error) {
	version := internal.TrimTagPrefix(release)
	artifactPrefix := fmt.Sprintf("%s_%s_", p.RepositoryName(), version)

	logger := p.Logger.With(slog.String("release", release))

	urlPrefix := fmt.Sprintf(p.RepositoryURL()+"/releases/download/%s/%s", release, artifactPrefix)

	v := Version{
		Version:             version,
		SHASumsURL:          urlPrefix + "SHA256SUMS",
		SHASumsSignatureURL: urlPrefix + "SHA256SUMS.sig",
	}

	checksums, err := p.GetSHASums(v.SHASumsURL)
	if err != nil {
		return nil, err
	}

	if checksums == nil {
		logger.Warn(fmt.Sprintf("checksums not found in release %s, skipping...", version))
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
			target.DownloadURL = fmt.Sprintf("%s/releases/download/v%s/%s", p.RepositoryURL(), version, target.Filename)
			target.SHASum, ok = checksums[target.Filename]
			if ok {
				v.Targets = append(v.Targets, target)
			}
		}
	}

	if len(v.Targets) == 0 {
		logger.Info("No artifacts in release, skipping...", slog.String("release", version))
		return nil, nil
	}

	v.Protocols, err = p.GetProtocols(urlPrefix + "manifest.json")
	if err != nil {
		return nil, err
	}
	if v.Protocols == nil {
		logger.Warn("Could not find manifest file, using default protocols")
		v.Protocols = []string{"5.0"}
	}

	return &v, nil
}
