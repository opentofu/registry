package github

import (
	"bytes"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

type Tag struct {
	Name   string
	Commit string
}

func parseTagsFromStdout(lines []string) ([]Tag, error) {
	tags := make([]Tag, 0, len(lines))

	for _, line := range lines {
		if !strings.Contains(line, "refs/tags/") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid format for tag '%s', expected two fields", line)
		}
		commit := fields[0]

		ref := fields[1]
		if !strings.HasPrefix(ref, "refs/tags/") {
			return nil, fmt.Errorf("invalid format for tag '%s', expected 'refs/tags/' prefix", line)
		}

		tag := strings.TrimPrefix(ref, "refs/tags/")
		if tag == "" {
			return nil, fmt.Errorf("invalid format for tag '%s', no version provided", line)
		}

		tags = append(tags, Tag{Name: tag, Commit: commit})
	}

	return tags, nil
}

// GetTags lists the tags of the remote repository and returns the refs/tags/ found
// Order is not guaranteed
func (c Client) GetTags(repositoryUrl string) ([]Tag, error) {
	done := c.cliThrottle()
	defer done()

	c.log.Info("Getting tags for repository", slog.String("repository", repositoryUrl))

	var buf bytes.Buffer
	var bufErr bytes.Buffer
	cmd := exec.Command("git", "ls-remote", "--tags", "--refs", repositoryUrl)
	cmd.Stdout = &buf
	cmd.Stderr = &bufErr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("could not get tags for %s, %w: %s", repositoryUrl, err, bufErr.String())
	}

	tags, err := parseTagsFromStdout(strings.Split(buf.String(), "\n"))
	if err != nil {
		return nil, fmt.Errorf("could not parse tags for %s: %w", repositoryUrl, err)
	}

	c.log.Info("Found tags for repository", slog.String("repository", repositoryUrl), slog.Int("count", len(tags)))
	return tags, nil
}
