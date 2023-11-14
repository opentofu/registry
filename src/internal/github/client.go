package github

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func EnvAuthToken() (string, error) {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		return "", fmt.Errorf("Expected $GH_TOKEN to be set, unable to authenticate with GitHub")
	}
	return token, nil

}

func NewGitHubClient(ctx context.Context, token string) *githubv4.Client {
	return githubv4.NewClient(getGithubOauth2Client(ctx, token))
}

func getGithubOauth2Client(ctx context.Context, token string) *http.Client {
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
}

func GetHttpRetryClient(ctx context.Context, token string) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Logger = nil
	retryClient.HTTPClient = getGithubOauth2Client(ctx, token)

	return retryClient.StandardClient()
}
