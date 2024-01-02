package github

import (
	"fmt"
	"log/slog"
)

type Repository struct {
	client Client
	log    *slog.Logger
	Owner  string
	Name   string
}

// URL constructs the URL to the repository on github.com.
func (r Repository) URL() string {
	return fmt.Sprintf("https://github.com/%s/%s", r.Owner, r.Name)
}

// VersionDownloadURL returns the location to download the repository from.
// git::https://github.com/terraform-aws-modules/terraform-aws-iam?ref=v5.30.0
func (r Repository) DownloadURL(tag string) string {
	return fmt.Sprintf("git::%s?ref=%s", r.URL(), tag)
}
