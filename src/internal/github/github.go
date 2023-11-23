package github

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/shurcooL/githubv4"
)

func (c Client) FetchPublishedReleases(owner string, repoName string) (releases []GHRelease, err error) {
	variables := map[string]interface{}{
		"owner":     githubv4.String(owner),
		"name":      githubv4.String(repoName),
		"perPage":   githubv4.Int(100),
		"endCursor": (*githubv4.String)(nil),
	}
	done := c.apiThrottle()
	defer done()

	for {
		var query GHRepository
		if err := c.ghClient.Query(c.ctx, &query, variables); err != nil {
			return nil, fmt.Errorf("failed to fetch releases for %s/%s: %w", owner, repoName, err)
		}

		c.log.Info("Checking for possible new releases", slog.Int("releases", len(query.Repository.Releases.Nodes)))

		for _, r := range query.Repository.Releases.Nodes {
			if r.IsDraft || r.IsPrerelease {
				continue
			}

			c.log.Info("New release fetched", slog.String("release", r.TagName), slog.String("created", r.CreatedAt.String()))
			releases = append(releases, r)
		}

		if !query.Repository.Releases.PageInfo.HasNextPage {
			c.log.Info("No more releases to fetch")
			break
		}

		variables["endCursor"] = githubv4.String(query.Repository.Releases.PageInfo.EndCursor)
	}

	c.log.Info("New releases fetched", slog.Int("releases", len(releases)))
	return releases, nil
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
