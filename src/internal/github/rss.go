package github

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/mmcdole/gofeed"
)

// GetTagsFromRss gets all tags found in the RSS feed of a GitHub releases page
// Tags are sorted by descending creation date
func (c Client) GetTagsFromRss(releasesRssUrl string) ([]string, error) {
	feed, err := c.getReleaseRssFeed(releasesRssUrl)
	if err != nil {
		return nil, err
	}

	var tags = make([]string, 0)
	for _, item := range feed.Items {
		tag, err := c.extractTag(item)
		if err != nil {
			return nil, err
		}
		if tag != nil {
			tags = append(tags, *tag)
		}
	}

	return tags, nil
}

func (c Client) extractTag(item *gofeed.Item) (*string, error) {
	pattern := regexp.MustCompile(`.*/(?P<Version>[a-zA-Z0-9.\-_+]+)$`)
	matches := pattern.FindStringSubmatch(item.GUID)

	if matches == nil {
		c.log.Warn(fmt.Sprintf("Could not parse RSS item %s", item.Link))
		return nil, nil
	}

	return &matches[pattern.SubexpIndex("Version")], nil
}

func (c Client) getReleaseRssFeed(releasesRssUrl string) (*gofeed.Feed, error) {
	done := c.rssThrottle()
	defer done()

	resp, err := c.httpClient.Get(releasesRssUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", releasesRssUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s got error %d", releasesRssUrl, resp.StatusCode)
	}

	return gofeed.NewParser().Parse(resp.Body)
}
