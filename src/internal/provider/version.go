package provider

import (
	"fmt"
	"log/slog"
	"registry-stable/internal"
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

func (p Provider) VersionFromTag(release string) (*Version, error) {
	version := internal.TrimTagPrefix(release)
	artifactPrefix := fmt.Sprintf("%s_%s_", p.RepositoryName(), version)

	urlPrefix := fmt.Sprintf(p.RepositoryURL()+"/releases/download/%s/%s", release, artifactPrefix)

	v := Version{
		Version:             version,
		SHASumsURL:          urlPrefix + "SHA256SUMS",
		SHASumsSignatureURL: urlPrefix + "SHA256SUMS.sig",
	}

	signatures, err := p.GetShaSums(v.SHASumsURL)
	if err != nil {
		return nil, err
	}
	if signatures == nil {
		p.Logger.Info("Signature not found in release, skipping...", slog.String("release", version))
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
			target.SHASum, ok = signatures[target.Filename]
			if ok {
				v.Targets = append(v.Targets, target)
			}
		}
	}

	if len(v.Targets) == 0 {
		p.Logger.Info("No artifacts in release, skipping...", slog.String("release", version))
		return nil, nil
	}

	v.Protocols, err = p.GetProtocols(urlPrefix + "manifest.json")
	if err != nil {
		return nil, err
	}
	if v.Protocols == nil {
		p.Logger.Warn("Could not find manifest file, using default protocols")
		v.Protocols = []string{"5.0"}
	}

	return &v, nil
}
