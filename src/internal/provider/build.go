package provider

import (
	"context"
	"fmt"
	"log"
	"slices"

	"registry-stable/internal"
	"registry-stable/internal/github"

	"golang.org/x/mod/semver"
)

func filterNewReleases(releases []github.GHRelease, existingMetadata MetadataFile) ([]github.GHRelease, error) {
	var existingVersions = make(map[string]bool)
	for _, v := range existingMetadata.Versions {
		existingVersions[v.Version] = true
	}

	var newReleases = make([]github.GHRelease, 0)
	for _, r := range releases {
		if !existingVersions[internal.TrimTagPrefix(r.TagName)] {
			newReleases = append(newReleases, r)
		}
	}

	log.Printf("Found %d releases that do not already exist in the metadata file", len(newReleases))

	return newReleases, nil
}

func BuildMetadataFile(p Provider, providerDataDir string) (*MetadataFile, error) {
	ctx := context.Background()

	token, err := github.EnvAuthToken()
	if err != nil {
		return nil, err
	}

	ghClient := github.NewGitHubClient(ctx, token)

	existingMetadata, err := p.ReadMetadata(providerDataDir)
	if err != nil {
		return nil, err
	}

	repoName := p.RepositoryName()
	releases, err := github.FetchPublishedReleases(ctx, ghClient, p.EffectiveNamespace(), repoName)
	if err != nil {
		return nil, err
	}

	newReleases, err := filterNewReleases(releases, existingMetadata)
	if err != nil {
		return nil, err
	}

	versions := make([]Version, 0)
	versionArtifactsMap := make(VersionArtifactsMap)

	for _, r := range newReleases {
		version := internal.TrimTagPrefix(r.TagName)
		versionArtifacts := getArtifacts(r)
		versionArtifactsMap[version] = versionArtifacts

		var targets = make([]Target, 0)
		for _, a := range versionArtifacts.TargetArtifacts {
			targets = append(targets, Target{
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

		versions = append(versions, Version{
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

	mergedMetadata := mergeMetadata(existingMetadata, MetadataFile{
		Versions: versions,
	})
	return &mergedMetadata, nil
}

func mergeMetadata(oldMetadata MetadataFile, newMetadata MetadataFile) MetadataFile {
	versions := append(newMetadata.Versions, oldMetadata.Versions...)

	semverSortFunc := func(a, b Version) int {
		return semver.Compare(fmt.Sprintf("s%s", a.Version), fmt.Sprintf("s%s", b.Version))
	}
	slices.SortFunc(versions, semverSortFunc)

	return MetadataFile{
		Repository: oldMetadata.Repository,
		Versions:   versions,
	}
}
