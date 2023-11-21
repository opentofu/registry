package github

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// TODO: probably move the Platform type inside providers as that's the only place this is used,
// and its a bit strange being in the github package

type Platform struct {
	OS   string
	Arch string
}

func DownloadAssetContents(ctx context.Context, logger *slog.Logger, downloadURL string) ([]byte, error) {
	logger = logger.With()

	token, err := EnvAuthToken()
	if err != nil {
		return nil, err
	}
	httpClient := GetHTTPRetryClient(token)

	logger.Info("Downloading asset", slog.String("url", downloadURL))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", downloadURL, err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error downloading asset %s: %w", downloadURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code when downloading asset %s: %d", downloadURL, resp.StatusCode)
	}

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset contents of %s: %w", downloadURL, err)
	}

	logger.Info("Asset downloaded", slog.String("url", downloadURL))

	return contents, nil
}
