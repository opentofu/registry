package github

import "fmt"

type Repository struct {
	client Client
	Owner  string
	Name   string
}

// URL constructs the URL to the repository on github.com.
func (r Repository) URL() string {
	return fmt.Sprintf("https://github.com/%s/%s", r.Owner, r.Name)
}

/*
func (m Module) TagsURL() string {
	repositoryUrl := m.RepositoryURL()
	return fmt.Sprintf("%s/tags.atom", repositoryUrl)
}

// RSSURL returns the URL of the RSS feed for the repository's releases.
func (p Provider) RSSURL() string {
	repositoryUrl := p.RepositoryURL()
	return fmt.Sprintf("%s/releases.atom", repositoryUrl)
}
*/

// VersionDownloadURL returns the location to download the repository from.
// The file should just contain a link to GitHub to download the tarball, ie:
// git::https://github.com/terraform-aws-modules/terraform-aws-iam?ref=v5.30.0
func (r Repository) DownloadURL(tag string) string {
	return fmt.Sprintf("git::%s?ref=%s", r.URL(), tag)
}

func (r Repository) ListTags() (Tags, error) {
	// TODO move Client to Repository
	return r.client.GetTags(r.URL())
}

func (r Repository) GetLatestReleases() (Tags, error) {
	// TODO move Client to Repository
	return r.client.GetTagsFromRSS(r.URL() + "/releases.atom")
}

func (r Repository) ReleaseAssetURL(release string, asset string) string {
	return fmt.Sprintf("%s/releases/download/%s/%s", r.URL(), release, asset)
}

func (r Repository) GetReleaseAsset(release string, asset string) ([]byte, error) {
	// TODO move Client to Repository
	return r.client.DownloadAssetContents(r.ReleaseAssetURL(release, asset))
}
