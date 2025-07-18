package github

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// DownloadAssetContents downloads the contents of the asset at the given URL and returns it directly
func (c Client) DownloadAssetContents(downloadURL string) ([]byte, error) {
	done := c.assetThrottle()
	defer done()

	logger := c.log.With(slog.String("url", downloadURL))
	logger.Info("Downloading asset")

	var resp *http.Response
	var err error
	maxRetries := 4

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = c.httpClient.Get(downloadURL)
		if err != nil {
			return nil, fmt.Errorf("error downloading asset %s: %w", downloadURL, err)
		}

		// Only retry on 403 Forbidden - suspected GitHub caching issue
		if resp.StatusCode == http.StatusForbidden && attempt < maxRetries-1 {
			resp.Body.Close()
			backoffDuration := time.Duration(1<<attempt) * time.Second

			logger.Warn("got 403 Forbidden, retrying after backoff due to suspected github caching issue 🤞",
				slog.Int("attempt", attempt+1),
				slog.Duration("backoff", backoffDuration))
			time.Sleep(backoffDuration)
			continue
		}

		break
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
