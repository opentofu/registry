package github

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/shurcooL/githubv4"
)

func EnvAuthToken() (string, error) {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		return "", fmt.Errorf("Expected $GH_TOKEN to be set, unable to authenticate with GitHub")
	}
	return token, nil

}

func NewGitHubClient(token string) *githubv4.Client {
	return githubv4.NewClient(getGithubOauth2Client(token))
}

type transport struct {
	token string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "OpenTofu Registry/1.0")
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}

func getGithubOauth2Client(token string) *http.Client {
	return &http.Client{Transport: &transport{token}}
}

func GetHTTPRetryClient(token string) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Logger = nil
	retryClient.HTTPClient = getGithubOauth2Client(token)

	return retryClient.StandardClient()
}
