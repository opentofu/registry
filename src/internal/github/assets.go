package github

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// DownloadAssetContents downloads the contents of the asset at the given URL and returns it directly
func (c Client) DownloadAssetContents(ctx context.Context, downloadURL string) ([]byte, error) {
	done := c.assetThrottle()
	defer done()

	logger := c.log.With(slog.String("url", downloadURL))
	logger.Info("Downloading asset")

	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error downloading asset %s: %w", downloadURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		logger.Warn("asset not found")
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code when downloading asset %s: %d", downloadURL, resp.StatusCode)
	}

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset contents of %s: %w", downloadURL, err)
	}

	logger.Info("asset successfully downloaded", slog.Int("size", len(contents)))

	return contents, nil
}
