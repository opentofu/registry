package github

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	httpInternal "registry-stable/internal/http"
)

type Platform struct {
	OS   string
	Arch string
}

func DownloadAssetContents(ctx context.Context, downloadURL string) (body io.ReadCloser, err error) {
	httpClient := httpInternal.GetHttpRetryClient()

	log.Printf("Downloading asset, url: %s", downloadURL)
	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if reqErr != nil {
		log.Printf("Failed to create request %s", reqErr)
		err = fmt.Errorf("failed to create request: %w", reqErr)
		return
	}

	resp, respErr := httpClient.Do(req)
	if respErr != nil {
		log.Printf("Error downloading asset %s", respErr)
		err = fmt.Errorf("error downloading asset: %w", respErr)
		return
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Printf("Unexpected status code when downloading asset: %d", resp.StatusCode)
		err = fmt.Errorf("unexpected status code when downloading asset: %d", resp.StatusCode)
		return
	}

	body = resp.Body

	log.Printf("Asset downloaded successfully")
	return
}

func ExtractPlatformFromFilename(filename string) *Platform {
	platformPattern := regexp.MustCompile(`.*_(?P<Os>[a-zA-Z0-9]+)_(?P<Arch>[a-zA-Z0-9]+).zip`)
	matches := platformPattern.FindStringSubmatch(filename)

	if matches == nil {
		return nil
	}

	platform := Platform{
		OS:   matches[platformPattern.SubexpIndex("Os")],
		Arch: matches[platformPattern.SubexpIndex("Arch")],
	}

	return &platform
}
