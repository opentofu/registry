package github

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/mmcdole/gofeed"
)

func (r Repository) GetLatestReleases() (Tags, error) {
	return r.GetTagsFromRSS("/releases.atom")
}

// GetTagsFromRSS gets all tags found in the RSS feed of a GitHub  page
// Tags are sorted by descending creation date
func (r Repository) GetTagsFromRSS(path string) (Tags, error) {
	feed, err := r.getRSSFeed(r.URL() + path)
	if err != nil {
		return nil, err
	}

	var tags = make(Tags, 0)
	for _, item := range feed.Items {
		tag := r.extractTag(item)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}

	return tags, nil
}

// tagPattern is used in extractTag to extract the tag from the RSS item
var tagPattern = regexp.MustCompile(`.*/(?P<Version>[a-zA-Z0-9.\-_+]+)$`)

func (r Repository) extractTag(item *gofeed.Item) *string {
	matches := tagPattern.FindStringSubmatch(item.GUID)

	if matches == nil {
		r.log.Warn(fmt.Sprintf("Could not parse RSS item %s", item.Link))
		return nil
	}

	return &matches[tagPattern.SubexpIndex("Version")]
}

func (r Repository) getRSSFeed(RSSURL string) (*gofeed.Feed, error) {
	done := r.client.rssThrottle()
	defer done()

	resp, err := r.client.httpClient.Get(RSSURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", RSSURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s got error %d", RSSURL, resp.StatusCode)
	}

	feed, err := gofeed.NewParser().Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to parse rss feed, %s: %w", RSSURL, err)
	}

	return feed, nil
}
