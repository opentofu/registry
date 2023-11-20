package provider

import (
	"context"
	"fmt"
	"registry-stable/internal/github"
	"strings"
	"sync"
)

type Artifact struct {
	Name        string
	DownloadURL string
}

type TargetArtifact struct {
	Artifact

	OS   string
	Arch string
}

type VersionArtifacts struct {
	ShaSumsArtifact          Artifact
	ShaSumsSignatureArtifact Artifact
	ManifestArtifact         Artifact
	TargetArtifacts          []TargetArtifact
}

type VersionArtifactsMap map[string]VersionArtifacts

type ShaSumResult struct {
	Index     int
	ShaSumMap map[string]string
	Err       error
}

type ProtocolsResult struct {
	Index     int
	Protocols []string
	Err       error
}

func getArtifacts(release github.GHRelease) VersionArtifacts {
	var shaSumsArtifact Artifact
	var shaSumsSignatureArtifact Artifact
	var manifestArtifact Artifact
	var targetArtifacts = make([]TargetArtifact, 0)

	for _, asset := range release.ReleaseAssets.Nodes {
		if platform := github.ExtractPlatformFromFilename(asset.Name); platform != nil {
			targetArtifacts = append(targetArtifacts, TargetArtifact{
				OS:   platform.OS,
				Arch: platform.Arch,
				Artifact: Artifact{
					Name:        asset.Name,
					DownloadURL: asset.DownloadURL,
				},
			})
		} else if strings.HasSuffix(asset.Name, "SHA256SUMS") {
			shaSumsArtifact = Artifact{
				Name:        asset.Name,
				DownloadURL: asset.DownloadURL,
			}
		} else if strings.HasSuffix(asset.Name, "SHA256SUMS.sig") {
			shaSumsSignatureArtifact = Artifact{
				Name:        asset.Name,
				DownloadURL: asset.DownloadURL,
			}
		} else if strings.HasSuffix(asset.Name, "_manifest.json") {
			manifestArtifact = Artifact{
				Name:        asset.Name,
				DownloadURL: asset.DownloadURL,
			}
		}
	}

	return VersionArtifacts{
		ShaSumsArtifact:          shaSumsArtifact,
		ShaSumsSignatureArtifact: shaSumsSignatureArtifact,
		ManifestArtifact:         manifestArtifact,
		TargetArtifacts:          targetArtifacts,
	}
}

// enrichWithDataFromArtifacts performs all data enrichment necessary for the provider, that require artifact downloads
// All necessary artifacts are downloaded in parallel, when calculating the SHAs and Protocols of each provider version
func enrichWithDataFromArtifacts(ctx context.Context, versions []Version, artifactsMap VersionArtifactsMap) ([]Version, error) {
	versionsCopy := versions

	shaSumCh := make(chan ShaSumResult, len(versionsCopy))
	protocolsCh := make(chan ProtocolsResult, len(versionsCopy))

	var wg sync.WaitGroup
	for i, v := range versionsCopy {
		wg.Add(2)

		go func(v Version, i int) {
			defer wg.Done()
			shaMap, err := GetShaSums(ctx, artifactsMap[v.Version].ShaSumsArtifact.DownloadURL)
			shaSumCh <- ShaSumResult{
				Index:     i,
				ShaSumMap: shaMap,
				Err:       err,
			}
		}(v, i)

		go func(v Version, i int) {
			defer wg.Done()
			protocols, err := GetProtocols(ctx, artifactsMap[v.Version].ManifestArtifact.DownloadURL)
			protocolsCh <- ProtocolsResult{
				Index:     i,
				Protocols: protocols,
				Err:       err,
			}
		}(v, i)
	}

	wg.Wait()
	close(shaSumCh)
	close(protocolsCh)

	for sr := range shaSumCh {
		if sr.Err != nil {
			return nil, fmt.Errorf("failed to find SHA of artifact: %w", sr.Err)
		}

		for i, t := range versionsCopy[sr.Index].Targets {
			if shaSum, ok := sr.ShaSumMap[t.Filename]; !ok {
				return nil, fmt.Errorf("failed to find SHA of file %s", t.Filename)
			} else {
				versionsCopy[sr.Index].Targets[i].SHASum = shaSum
			}
		}
	}

	for pr := range protocolsCh {
		if pr.Err != nil {
			return nil, fmt.Errorf("failed to get supported protocols for provider version: %w", pr.Err)
		}

		versionsCopy[pr.Index].Protocols = pr.Protocols
	}

	return versionsCopy, nil
}
