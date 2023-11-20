package provider

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"registry-stable/internal"
	"registry-stable/internal/github"
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

func (p Provider) VersionFromRelease(release github.GHRelease) (*Version, error) {
	version := internal.TrimTagPrefix(release.TagName)
	artifactPrefix := fmt.Sprintf("%s_%s", p.RepositoryName(), version)

	assets := make(map[string]string)
	for _, asset := range release.ReleaseAssets.Nodes {
		assets[asset.Name] = asset.DownloadURL
	}

	var ok bool
	ctx := context.TODO()
	v := Version{Version: version}

	for _, os := range goos {
		for _, arch := range goarch {
			target := Target{
				OS:       os,
				Arch:     arch,
				Filename: fmt.Sprintf("%s_%s_%s.zip", artifactPrefix, os, arch),
			}
			target.DownloadURL, ok = assets[target.Filename]
			if !ok {
				// Artifact does not exist for this platform
				continue
			}

			v.Targets = append(v.Targets, target)
		}
	}

	if len(v.Targets) == 0 {
		log.Printf("could not find artifacts in release of provider %s version %s, skipping...", p.ProviderName, version)
		return nil, nil
	}

	v.SHASumsURL, ok = assets[fmt.Sprintf("%s_%s", artifactPrefix, "SHA256SUMS")]
	if !ok {
		return nil, fmt.Errorf("Provider %s release %s missing SHA256SUMS", p.RepositoryName(), version)
	}

	v.SHASumsSignatureURL, ok = assets[fmt.Sprintf("%s_%s", artifactPrefix, "SHA256SUMS.sig")]
	if !ok {
		return nil, fmt.Errorf("Provider %s release %s missing SHA256SUMS.sig", p.RepositoryName(), version)
	}

	signatures, err := GetShaSums(ctx, v.SHASumsURL)
	if err != nil {
		return nil, err
	}

	for i, target := range v.Targets {
		target.SHASum, ok = signatures[target.Filename]
		if !ok {
			return nil, fmt.Errorf("Provider %s release %s missing signature for %s", p.RepositoryName(), version, target.Filename)
		}
		v.Targets[i] = target
	}

	manifestUrl, ok := assets[fmt.Sprintf("%s_%s", artifactPrefix, "manifest.json")]
	if !ok {
		slog.Warn("Could not find manifest file, using default protocols")
		v.Protocols = []string{"5.0"}
	} else {
		v.Protocols, err = GetProtocols(ctx, manifestUrl)
		if err != nil {
			return nil, err
		}
	}

	return &v, nil
}
