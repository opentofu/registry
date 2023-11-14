package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"registry-stable/internal"
	"registry-stable/internal/github"
	"registry-stable/internal/provider"
	"strings"
)

func BuildMetadataFile(p provider.Provider) (*provider.MetadataFile, error) {
	ctx := context.Background()
	ghClient := github.NewGitHubClient(ctx, os.Getenv("GH_TOKEN"))

	repoName := p.RepositoryName()
	releases, err := github.FetchPublishedReleases(ctx, ghClient, p.EffectiveNamespace(), repoName)
	if err != nil {
		return nil, err
	}

	versions := make([]provider.Version, 0)
	for _, r := range releases {
		var shaSumsArtifact github.ReleaseAsset
		var shaSumsSignatureArtifact github.ReleaseAsset

		var targets = make([]provider.Target, 0)
		for _, asset := range r.ReleaseAssets.Nodes {
			if platform := github.ExtractPlatformFromFilename(asset.Name); platform != nil {
				if err != nil {
					return nil, err
				}
				targets = append(targets, provider.Target{
					OS:          platform.OS,
					Arch:        platform.Arch,
					Filename:    asset.Name,
					DownloadURL: asset.DownloadURL,
				})
			} else if strings.HasSuffix(asset.Name, "SHA256SUMS") {
				shaSumsArtifact = asset
			} else if strings.HasSuffix(asset.Name, "SHA256SUMS.sig") {
				shaSumsSignatureArtifact = asset
			}
		}
		if len(targets) == 0 {
			log.Printf("could not find artifacts in release of provider %s version %s, skipping...", p.ProviderName, r.TagName)
			continue
		}
		if (shaSumsArtifact == github.ReleaseAsset{}) {
			return nil, fmt.Errorf("could not SHASUMS artifact for provider %s version %s", p.ProviderName, r.TagName)
		}
		if (shaSumsSignatureArtifact == github.ReleaseAsset{}) {
			return nil, fmt.Errorf("could not SHASUMS signature artifact for provider %s version %s", p.ProviderName, r.TagName)
		}

		versions = append(versions, provider.Version{
			Version:             internal.TrimTagPrefix(r.TagName),
			Protocols:           []string{"5.0"},
			SHASumsURL:          shaSumsArtifact.DownloadURL,
			SHASumsSignatureURL: shaSumsSignatureArtifact.DownloadURL,
			Targets:             targets,
		})
	}

	// TODO all asset downloads - Shasums and figuring out the protocols
	//versions, err = enrichWithShaSums(ctx, versions)
	//if err != nil {
	//	return nil, err
	//}

	return &provider.MetadataFile{
		Versions: versions,
	}, nil

}
