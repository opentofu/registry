package github

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

type Tag struct {
	Ref    string
	Commit string
}

func parseTagsFromStdout(lines []string) ([]Tag, error) {
	tags := make([]Tag, 0, len(lines))

	for _, line := range lines {
		prefix := "refs/tags/"
		if !strings.Contains(line, prefix) {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid format for tag %q, expected two fields", line)
		}

		commit := fields[0]
		if _, err := hex.DecodeString(commit); err != nil {
			return nil, fmt.Errorf("invalid format for commit %q: %w", line, err)
		}

		ref := fields[1]
		if !strings.HasPrefix(ref, prefix) {
			return nil, fmt.Errorf("invalid format for tag %q, expected %q prefix", line, prefix)
		}

		tag := strings.TrimPrefix(ref, prefix)
		if tag == "" {
			return nil, fmt.Errorf("invalid format for tag %q, no version provided", line)
		}

		tags = append(tags, Tag{Ref: tag, Commit: commit})
	}

	return tags, nil
}

// GetTags lists the tags of the remote repository and returns the refs/tags/ found
func (c Client) GetTags(repositoryURL string) ([]Tag, error) {
	done := c.cliThrottle()
	defer done()

	c.log.Info("Getting tags for repository", slog.String("repository", repositoryURL))

	var buf bytes.Buffer
	var bufErr bytes.Buffer
	cmd := exec.Command("git", "ls-remote", "--tags", "--refs", repositoryURL)
	cmd.Stdout = &buf
	cmd.Stderr = &bufErr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("could not get tags for %s, %w: %s", repositoryURL, err, bufErr.String())
	}

	tags, err := parseTagsFromStdout(strings.Split(buf.String(), "\n"))
	if err != nil {
		return nil, fmt.Errorf("could not parse tags for %s: %w", repositoryURL, err)
	}

	c.log.Info("Found tags for repository", slog.String("repository", repositoryURL), slog.Int("count", len(tags)))
	return tags, nil
}
