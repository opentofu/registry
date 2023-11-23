package github

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/time/rate"
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
	cliLimiter *rate.Limiter
}

func NewClient(ctx context.Context, log *slog.Logger, token string) Client {
	// These rate limits are guesses
	httpClient := &http.Client{Transport: &transport{
		token:       token,
		ctx:         ctx,
		ratelimiter: rate.NewLimiter(rate.Every(time.Second), 1),
	}}
	httpClientFast := &http.Client{Transport: &transport{
		ctx:         ctx,
		ratelimiter: rate.NewLimiter(rate.Every(time.Second/30), 1),
	}}

	return Client{
		ctx:        ctx,
		log:        log.WithGroup("github"),
		httpClient: httpClientFast,
		ghClient:   githubv4.NewClient(httpClient),
		cliLimiter: rate.NewLimiter(rate.Every(time.Second/10), 1),
	}
	/* TODO
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Logger = nil
	retryClient.HTTPClient = getGithubOauth2Client(token)
	*/
}

func (c Client) WithLogger(log *slog.Logger) Client {
	return Client{
		ctx:        c.ctx,
		log:        log.WithGroup("github"),
		httpClient: c.httpClient,
		ghClient:   c.ghClient,
		cliLimiter: c.cliLimiter,
	}
}

type transport struct {
	token       string
	ctx         context.Context
	ratelimiter *rate.Limiter
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	err := t.ratelimiter.Wait(t.ctx)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(t.ctx)
	req.Header.Set("User-Agent", "OpenTofu Registry/1.0")
	req.Header.Set("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}
