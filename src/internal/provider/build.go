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

func (existingMetadata MetadataFile) filterNewReleases(releases []github.GHRelease) []github.GHRelease {
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

	return newReleases
}

func (p Provider) buildMetadataFile() (*MetadataFile, error) {
	ctx := context.Background()

	token, err := github.EnvAuthToken()
	if err != nil {
		return nil, err
	}

	ghClient := github.NewGitHubClient(ctx, token)

	meta, err := p.ReadMetadata()
	if err != nil {
		return nil, err
	}

	releases, err := github.FetchPublishedReleases(ctx, ghClient, p.EffectiveNamespace(), p.RepositoryName())
	if err != nil {
		return nil, err
	}

	releases = meta.filterNewReleases(releases)

	versions := make([]Version, 0)
	versionArtifactsMap := make(VersionArtifactsMap)

	for _, r := range releases {
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

	meta.Versions = append(meta.Versions, versions...)

	semverSortFunc := func(a, b Version) int {
		return semver.Compare(fmt.Sprintf("s%s", a.Version), fmt.Sprintf("s%s", b.Version))
	}
	slices.SortFunc(meta.Versions, semverSortFunc)

	return &meta, nil
}
