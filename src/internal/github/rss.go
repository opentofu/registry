package github

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/mmcdole/gofeed"
)

// GetTagsFromRSS gets all tags found in the RSS feed of a GitHub releases page
// Tags are sorted by descending creation date
func (c Client) GetTagsFromRSS(releasesRSSURL string) ([]string, error) {
	feed, err := c.getReleaseRSSFeed(releasesRSSURL)
	if err != nil {
		return nil, err
	}

	var tags = make([]string, 0)
	for _, item := range feed.Items {
		tag := c.extractTag(item)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}

	return tags, nil
}

// tagPattern is used in extractTag to extract the tag from the RSS item
var tagPattern = regexp.MustCompile(`.*/(?P<Version>[a-zA-Z0-9.\-_+]+)$`)

func (c Client) extractTag(item *gofeed.Item) *string {
	matches := tagPattern.FindStringSubmatch(item.GUID)

	if matches == nil {
		c.log.Warn(fmt.Sprintf("Could not parse RSS item %s", item.Link))
		return nil
	}

	return &matches[tagPattern.SubexpIndex("Version")]
}

func (c Client) getReleaseRSSFeed(releasesRSSURL string) (*gofeed.Feed, error) {
	done := c.rssThrottle()
	defer done()

	resp, err := c.httpClient.Get(releasesRSSURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", releasesRSSURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s got error %d", releasesRSSURL, resp.StatusCode)
	}

	feed, err := gofeed.NewParser().Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to parse rss feed, %s: %w", releasesRSSURL, err)
	}

	return feed, nil
}
