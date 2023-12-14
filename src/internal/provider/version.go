package provider

import (
	"fmt"
	"log/slog"

	"github.com/opentofu/registry-stable/internal"
	"github.com/opentofu/registry-stable/internal/re"
)

// VersionFromTag fetches information about an individual release based on the GitHub release name
func (p Provider) VersionFromTag(release string) (*Version, error) {
	version := internal.TrimTagPrefix(release)

	logger := p.Log.With(slog.String("release", release))

	versionAssetPrefix := fmt.Sprintf("%s_%s_", p.Repository.Name, version)
	shasumName := versionAssetPrefix + "SHA256SUMS"
	shasumSig := versionAssetPrefix + "SHA256SUMS.sig"

	v := Version{
		Version:             version,
		SHASumsURL:          p.Repository.ReleaseAssetURL(release, shasumName),
		SHASumsSignatureURL: p.Repository.ReleaseAssetURL(release, shasumSig),
	}

	checkdata, err := p.Repository.GetReleaseAsset(release, shasumName)
	if err != nil {
		return nil, err
	}
	if checkdata == nil {
		logger.Warn("checksums not found in release, skipping...")
		return nil, nil
	}

	checksums := shaFileToMap(checkdata)

	// repo_v?version_(os)_(arch).zip
	releaseAssetMatcher := re.MustCompile(fmt.Sprintf("%s_v*%s_(?P<OS>\\w+)_(?P<ARCH>\\w+).zip", p.Repository.Name, version))

	for filename, sum := range checksums {
		match := releaseAssetMatcher.Match(filename)
		if match == nil {
			logger.Warn("Invalid file in release", slog.String("asset", filename))
			continue
		}
		v.Targets = append(v.Targets, Target{
			OS:          match["OS"],
			Arch:        match["ARCH"],
			Filename:    filename,
			SHASum:      sum,
			DownloadURL: p.Repository.ReleaseAssetURL(release, filename),
		})
	}

	if len(v.Targets) == 0 {
		logger.Info("No artifacts in release, skipping...", slog.String("release", version))
		return nil, nil
	}

	manifestData, err := p.Repository.GetReleaseAsset(release, versionAssetPrefix+"manifest.json")
	if err != nil {
		return nil, err
	}
	if manifestData != nil {
		manifest, err := parseManifestContents(manifestData)
		if err != nil {
			logger.Warn("Manifest file invalid, ignoring...", slog.Any("err", err))
		} else {
			v.Protocols = manifest.Metadata.ProtocolVersions
		}
	}

	if len(v.Protocols) == 0 {
		logger.Warn("Could not find protocols in manifest file, using default protocols")
		// TODO move this to the generator
		v.Protocols = []string{"5.0"}
	}

	return &v, nil
}
