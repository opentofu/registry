package github

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/shurcooL/githubv4"
)

func EnvAuthToken() (string, error) {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		return "", fmt.Errorf("Expected $GH_TOKEN to be set, unable to authenticate with GitHub")
	}
	return token, nil

}

type Client struct {
	ctx        context.Context
	log        *slog.Logger
	httpClient *http.Client
	ghClient   *githubv4.Client
}

func NewClient(ctx context.Context, log *slog.Logger, token string) Client {
	httpClient := &http.Client{Transport: &transport{token}}

	return Client{
		ctx:        ctx,
		log:        log.WithGroup("github"),
		httpClient: httpClient,
		ghClient:   githubv4.NewClient(httpClient),
	}
	/* TODO
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Logger = nil
	retryClient.HTTPClient = getGithubOauth2Client(token)
	*/
	// TODO rate limiting
}

func (c Client) WithLogger(log *slog.Logger) Client {
	return Client{
		ctx:        c.ctx,
		log:        log.WithGroup("github"),
		httpClient: c.httpClient,
		ghClient:   c.ghClient,
	}
}

type transport struct {
	token string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "OpenTofu Registry/1.0")
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}
