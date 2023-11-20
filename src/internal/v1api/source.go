package v1api

import (
	"path/filepath"
	"registry-stable/internal/github"
	"registry-stable/internal/provider"
)

type ProviderSource struct {
	Provider provider.Provider
	Meta     provider.MetadataFile
}

func (p ProviderSource) VersionListingPath() string {
	return filepath.Join("v1", "providers", p.Provider.Namespace, p.Provider.ProviderName, "versions")
}

func (p ProviderSource) VersionDownloadPath(ver provider.Version, details ProviderVersionDetails) string {
	return filepath.Join("v1", "providers", p.Provider.Namespace, p.Provider.ProviderName, ver.Version, "download", details.OS, details.Arch)
}

func (p ProviderSource) Versions() ProviderVersionListingResponse {
	versions := make([]ProviderVersionResponseItem, len(p.Meta.Versions))

	for versionIdx, ver := range p.Meta.Versions {
		verResp := ProviderVersionResponseItem{
			Version:   ver.Version,
			Protocols: ver.Protocols,
			Platforms: make([]github.Platform, len(ver.Targets)),
		}

		for targetIdx, target := range ver.Targets {
			verResp.Platforms[targetIdx] = github.Platform{
				OS:   target.OS,
				Arch: target.Arch,
			}
		}
		versions[versionIdx] = verResp
	}

	return ProviderVersionListingResponse{versions}
}

func (p ProviderSource) VersionDetails() map[string]ProviderVersionDetails {
	versionDetails := make(map[string]ProviderVersionDetails)

	for _, ver := range p.Meta.Versions {
		for _, target := range ver.Targets {
			details := ProviderVersionDetails{
				Protocols:           ver.Protocols,
				OS:                  target.OS,
				Arch:                target.Arch,
				Filename:            target.Filename,
				DownloadURL:         target.DownloadURL,
				SHASumsURL:          ver.SHASumsURL,
				SHASumsSignatureURL: ver.SHASumsSignatureURL,
				SHASum:              target.SHASum,
				SigningKeys:         SigningKeys{}, // TODO: Add gpg keys
			}
			versionDetails[p.VersionDownloadPath(ver, details)] = details
		}
	}
	return versionDetails
}
