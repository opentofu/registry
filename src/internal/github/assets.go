package github

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func (r Repository) ReleaseAssetURL(release string, asset string) string {
	return fmt.Sprintf("%s/releases/download/%s/%s", r.URL(), release, asset)
}

// GetReleaseAsset downloads the contents of the asset and returns it directly
func (r Repository) GetReleaseAsset(release string, asset string) ([]byte, error) {
	done := r.client.assetThrottle()
	defer done()

	downloadURL := r.ReleaseAssetURL(release, asset)

	logger := r.log.With(slog.String("url", downloadURL))
	logger.Info("Downloading asset")

	resp, err := r.client.httpClient.Get(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading asset %s: %w", downloadURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// TODO specific error instead of nil
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
