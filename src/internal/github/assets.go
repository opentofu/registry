package github

import (
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

func (c Client) DownloadAssetContents(downloadURL string) ([]byte, error) {
	done := c.assetThrottle()
	defer done()

	c.log.Info("Downloading asset", slog.String("url", downloadURL))

	resp, err := c.httpClient.Get(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading asset %s: %w", downloadURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		c.log.Info("Asset not found", slog.String("url", downloadURL))
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code when downloading asset %s: %d", downloadURL, resp.StatusCode)
	}

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset contents of %s: %w", downloadURL, err)
	}

	c.log.Info("Asset downloaded", slog.String("url", downloadURL))

	return contents, nil
}
