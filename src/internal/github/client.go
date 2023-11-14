package github

import (
	"context"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"net/http"
)

func NewGitHubClient(ctx context.Context, token string) *githubv4.Client {
	return githubv4.NewClient(getGithubOauth2Client(ctx, token))
}

func getGithubOauth2Client(ctx context.Context, token string) *http.Client {
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
}
