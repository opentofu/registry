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

func (p Provider) ArtifactName(version string, suffix string) string {
	return fmt.Sprintf("%s_%s_%s", p.RepositoryName(), version, suffix)
}

func (p Provider) ArtifactURL(release string, version string, suffix string) string {
	return fmt.Sprintf("%s/releases/download/%s/%s", p.RepositoryURL(), release, p.ArtifactName(version, suffix))
}

// VersionFromTag fetches information about an individual release based on the GitHub release name
func (p Provider) VersionFromTag(release string) (*Version, error) {
	version := internal.TrimTagPrefix(release)

	logger := p.Logger.With(slog.String("release", release))

	v := Version{
		Version:             version,
		SHASumsURL:          p.ArtifactURL(release, version, "SHA256SUMS"),
		SHASumsSignatureURL: p.ArtifactURL(release, version, "SHA256SUMS.sig"),
	}

	checksums, err := p.GetSHASums(v.SHASumsURL)
	if err != nil {
		return nil, err
	}

	if checksums == nil {
		// Attempt to use workaround
		logger.Warn("Attempting shasum workaround for release/version not using the same prefix")
		v.SHASumsURL = p.ArtifactURL(release, "v"+version, "SHA256SUMS")
		v.SHASumsSignatureURL = p.ArtifactURL(release, "v"+version, "SHA256SUMS.sig")
		checksums, err = p.GetSHASums(v.SHASumsURL)
		if err != nil {
			return nil, err
		}
	}

	if checksums == nil {
		logger.Warn("checksums not found in release, skipping...")
		return nil, nil
	}

	var ok bool
	for _, os := range goos {
		for _, arch := range goarch {
			suffix := fmt.Sprintf("%s_%s.zip", os, arch)

			target := Target{
				OS:          os,
				Arch:        arch,
				Filename:    p.ArtifactName(version, suffix),
				DownloadURL: p.ArtifactURL(release, version, suffix),
			}
			target.SHASum, ok = checksums[target.Filename]
			if !ok {
				// now try and pull it with the v in the version
				target.Filename = p.ArtifactName("v"+version, suffix)
				target.DownloadURL = p.ArtifactURL(release, "v"+version, suffix)
				target.SHASum, ok = checksums[target.Filename]
				if !ok {
					// Release target without checksum (invalid)
					continue
				}
			}

			target.Hash1, err = p.CalculateHash1(target.DownloadURL, target.SHASum)
			if err != nil {
				// Release is inaccessable, misconfigured, or corrupt
				return nil, err
			}

			v.Targets = append(v.Targets, target)
		}
	}

	if len(v.Targets) == 0 {
		logger.Info("No artifacts in release, skipping...", slog.String("release", version))
		return nil, nil
	}

	v.Protocols, err = p.GetProtocols(p.ArtifactURL(release, version, "manifest.json"))
	if err != nil {
		return nil, err
	}
	if v.Protocols == nil {
		logger.Warn("Attempting protocol workaround for release/version not using the same prefix")
		v.Protocols, err = p.GetProtocols(p.ArtifactURL(release, "v"+version, "manifest.json"))
		if err != nil {
			return nil, err
		}
	}

	if v.Protocols == nil {
		logger.Warn("Could not find manifest file, using default protocols")
		v.Protocols = []string{"5.0"}
	}

	return &v, nil
}
