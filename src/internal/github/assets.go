package github

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

// TODO: probably move the Platform type inside providers as that's the only place this is used,
// and its a bit strange being in the github package

type Platform struct {
	OS   string
	Arch string
}

func DownloadAssetContents(ctx context.Context, downloadURL string) ([]byte, error) {
	token, err := EnvAuthToken()
	if err != nil {
		return nil, err
	}
	httpClient := GetHTTPRetryClient(token)

	log.Printf("Downloading asset, url: %s", downloadURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		log.Printf("Failed to create request %s", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error downloading asset %s", err)
		return nil, fmt.Errorf("error downloading asset: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Printf("Unexpected status code when downloading asset: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code when downloading asset: %d", resp.StatusCode)
	}

	log.Printf("Asset downloaded successfully")

	contents, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read asset contents: %w", err)
	}

	return contents, nil
}
