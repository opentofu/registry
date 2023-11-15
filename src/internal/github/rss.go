package github

import (
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"net/http"
	"os"
	"regexp"
	httpInternal "registry-stable/internal/http"
)

// GetTagsFromRss gets all tags found in the RSS feed of a GitHub releases page
// Tags are sorted by descending creation date
func GetTagsFromRss(releasesRssUrl string) ([]string, error) {
	feed, err := getReleaseRssFeed(releasesRssUrl)
	if err != nil {
		return nil, err
	}

	var tags = make([]string, 0)
	for _, item := range feed.Items {
		tag, err := extractTag(item)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func extractTag(item *gofeed.Item) (string, error) {
	pattern := regexp.MustCompile(`.*/(?P<Version>[a-zA-Z0-9.\-_+]+)$`)
	matches := pattern.FindStringSubmatch(item.GUID)

	if matches == nil {
		return "", fmt.Errorf("could not parse RSS item %s", item.Link)
	}

	return matches[pattern.SubexpIndex("Version")], nil
}

func getReleaseRssFeed(releasesRssUrl string) (feed *gofeed.Feed, err error) {
	client := httpInternal.GetHttpRetryClient()

	req, err := http.NewRequestWithContext(context.Background(), "GET", releasesRssUrl, nil)
	if err != nil {
		return nil, err
	}

	// TODO Commonize?
	token := os.Getenv("GH_TOKEN")

	req.Header.Set("User-Agent", "OpenTofu Registry/1.0")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			ce := resp.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("got error %d", resp.StatusCode)
	}

	return gofeed.NewParser().Parse(resp.Body)
}
