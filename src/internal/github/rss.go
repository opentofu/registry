package github

import (
	"context"
	"fmt"
	"regexp"

	"github.com/mmcdole/gofeed"
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
	token, err := EnvAuthToken()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client := GetHttpRetryClient(ctx, token)

	resp, err := client.Get(releasesRssUrl)
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
