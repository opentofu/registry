package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"registry-stable/internal"
	"registry-stable/internal/github"
	"registry-stable/internal/provider"
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
	versionArtifactsMap := make(VersionArtifactsMap)

	for _, r := range releases {
		version := internal.TrimTagPrefix(r.TagName)
		versionArtifacts := getArtifacts(r)
		versionArtifactsMap[version] = versionArtifacts

		var targets = make([]provider.Target, 0)
		for _, a := range versionArtifacts.TargetArtifacts {
			targets = append(targets, provider.Target{
				OS:          a.OS,
				Arch:        a.Arch,
				Filename:    a.Name,
				DownloadURL: a.DownloadURL,
			})
		}
		if len(targets) == 0 {
			log.Printf("could not find artifacts in release of provider %s version %s, skipping...", p.ProviderName, r.TagName)
			continue
		}
		if (versionArtifacts.ShaSumsArtifact == Artifact{}) {
			return nil, fmt.Errorf("could not SHASUMS artifact for provider %s version %s", p.ProviderName, r.TagName)
		}
		if (versionArtifacts.ShaSumsSignatureArtifact == Artifact{}) {
			return nil, fmt.Errorf("could not SHASUMS signature artifact for provider %s version %s", p.ProviderName, r.TagName)
		}

		versions = append(versions, provider.Version{
			Version:             version,
			SHASumsURL:          versionArtifacts.ShaSumsArtifact.DownloadURL,
			SHASumsSignatureURL: versionArtifacts.ShaSumsSignatureArtifact.DownloadURL,
			Targets:             targets,
		})
	}

	versions, err = enrichWithDataFromArtifacts(ctx, versions, versionArtifactsMap)
	if err != nil {
		return nil, err
	}

	return &provider.MetadataFile{
		Versions: versions,
	}, nil

}
