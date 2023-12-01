package github

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/shurcooL/githubv4"
)

const UserAgent = "OpenTofu Registry/1.0"

// EnvAuthToken returns the GitHub token from the environment.
func EnvAuthToken() (string, error) {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		return "", fmt.Errorf("expected the GH_TOKEN environment variable to be set, unable to authenticate with GitHub")
	}
	return token, nil
}

// Client is a GitHub client that abstracts away the different GitHub APIs and handles rate limiting/throttling.
type Client struct {
	ctx        context.Context
	log        *slog.Logger
	httpClient *http.Client
	ghClient   *githubv4.Client

	cliThrottle   Throttle
	apiThrottle   Throttle
	assetThrottle Throttle
	rssThrottle   Throttle
}

// NewClient creates a new GitHub client.
func NewClient(ctx context.Context, log *slog.Logger, token string) Client {
	httpClient := &http.Client{Transport: &transport{token: token, ctx: ctx}}
	return Client{
		ctx:        ctx,
		log:        log.WithGroup("github"),
		httpClient: httpClient,
		ghClient:   githubv4.NewClient(httpClient),

		cliThrottle:   NewThrottle(ctx, time.Second/60, 60),
		apiThrottle:   NewThrottle(ctx, time.Second, 3),
		assetThrottle: NewThrottle(ctx, time.Second/60, 30),
		rssThrottle:   NewThrottle(ctx, time.Second/30, 30),
	}
	/* TODO
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Logger = nil
	retryClient.HTTPClient = getGithubOauth2Client(token)
	*/
}

// WithLogger returns a new Client with the given logger.
func (c Client) WithLogger(log *slog.Logger) Client {
	return Client{
		ctx:        c.ctx,
		log:        log.WithGroup("github"),
		httpClient: c.httpClient,
		ghClient:   c.ghClient,

		cliThrottle:   c.cliThrottle,
		apiThrottle:   c.apiThrottle,
		assetThrottle: c.assetThrottle,
		rssThrottle:   c.rssThrottle,
	}
}

// transport is a http.RoundTripper that makes sure all requests have the
// correct User-Agent and Authorization headers set.
type transport struct {
	token  string
	ctx    context.Context
	parent http.Transport
}

// RoundTrip is needed to implement the http.RoundTripper interface.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.WithContext(t.ctx)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Authorization", "Bearer "+t.token)

	parent := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return parent.RoundTrip(req)
}
