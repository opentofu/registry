package github

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/shurcooL/githubv4"
)

func FetchPublishedReleases(ctx context.Context, ghClient *githubv4.Client, owner string, repoName string) (releases []GHRelease, err error) {
	variables := map[string]interface{}{
		"owner":     githubv4.String(owner),
		"name":      githubv4.String(repoName),
		"perPage":   githubv4.Int(100),
		"endCursor": (*githubv4.String)(nil),
	}

	for {
		nodes, endCursor, fetchErr := fetchReleaseNodes(ctx, ghClient, variables)
		if fetchErr != nil {
			log.Printf("Failed to fetch release nodes")
			return nil, fmt.Errorf("failed to fetch release nodes: %w", fetchErr)
		}

		log.Printf("Checking for possible new releases: %d", len(nodes))

		for _, r := range nodes {
			if r.IsDraft || r.IsPrerelease {
				continue
			}

			log.Printf("New release fetched. Release: %s, Created at: %s", r.TagName, r.CreatedAt)
			releases = append(releases, r)
		}

		if endCursor == nil {
			log.Printf("No more releases to fetch")
			break
		}

		variables["endCursor"] = githubv4.String(*endCursor)
	}

	log.Printf("New releases fetched: %d", len(releases))
	return releases, nil
}

func fetchReleaseNodes(ctx context.Context, ghClient *githubv4.Client, variables map[string]interface{}) (releases []GHRelease, endCursor *string, err error) {
	var query GHRepository

	if queryErr := ghClient.Query(ctx, &query, variables); queryErr != nil {
		return nil, nil, fmt.Errorf("failed to query for releases: %w", queryErr)
	}

	if query.Repository.Releases.PageInfo.HasNextPage {
		endCursor = &query.Repository.Releases.PageInfo.EndCursor
	}

	releases = query.Repository.Releases.Nodes

	return releases, endCursor, err
}

// GHRelease represents a release on GitHub.
// This provides details about the release, including its tag name, release assets, and its release status (draft, prerelease, etc.).
type GHRelease struct {
	ID            string // The ID of the release.
	TagName       string // The tag name associated with the release.
	ReleaseAssets struct {
		Nodes []ReleaseAsset // A list of assets for the release.
	} `graphql:"releaseAssets(first:100)"`
	IsDraft      bool     // Indicates if the release is a draft.
	IsLatest     bool     // Indicates if the release is the latest.
	IsPrerelease bool     // Indicates if the release is a prerelease.
	TagCommit    struct { // The commit associated with the release tag.
		//nolint: revive, stylecheck // This is a struct provided by the GitHub GraphQL API.
		TarballUrl string // The URL to download the release tarball.
	}
	CreatedAt time.Time // The time the release was created.
}

// ReleaseAsset represents a single asset within a GitHub release.
// This includes details such as the download URL and the name of the asset.
type ReleaseAsset struct {
	ID          string // The ID of the asset.
	DownloadURL string // The URL to download the asset.
	Name        string // The name of the asset.
}

// GHRepository encapsulates GitHub repository details with a focus on its releases.
// This is structured to align with the expected response format from GitHub's GraphQL API.
type GHRepository struct {
	Repository struct {
		Releases struct {
			PageInfo struct {
				HasNextPage bool   // Indicates if there are more pages of releases.
				EndCursor   string // The cursor for pagination.
			}
			Nodes []GHRelease // A list of GitHub releases.
		} `graphql:"releases(first: $perPage, orderBy: {field: CREATED_AT, direction: DESC}, after: $endCursor)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}
