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
}*/

// VersionDownloadURL returns the location to download the repository from.
// The file should just contain a link to GitHub to download the tarball, ie:
// git::https://github.com/terraform-aws-modules/terraform-aws-iam?ref=v5.30.0
func (r Repository) DownloadURL(tag string) string {
	return fmt.Sprintf("git::%s?ref=%s", r.URL(), tag)
}

func (r Repository) ListTags() (Tags, error) {
	// TODO move GetTags from Client to Repository
	return r.client.GetTags(r.URL())
}
