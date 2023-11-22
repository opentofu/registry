package github

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func parseTagsFromStdout(lines []string) ([]string, error) {
	tags := make([]string, 0, len(lines))

	for _, line := range lines {
		if !strings.Contains(line, "refs/tags/") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("invalid format for tag '%s', expected two fields", line)
		}

		ref := fields[1]
		if !strings.HasPrefix(ref, "refs/tags/") {
			return nil, fmt.Errorf("invalid format for tag '%s', expected 'refs/tags/' prefix", line)
		}

		tag := strings.TrimPrefix(ref, "refs/tags/")
		if tag == "" {
			return nil, fmt.Errorf("invalid format for tag '%s', no version provided", line)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

// GetTags lists the tags of the remote repository and returns the refs/tags/ found
func (client Client) GetTags(repositoryUrl string) ([]string, error) {
	err := client.cliLimiter.Wait(client.ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("Getting tags for repository %s", repositoryUrl)

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

	log.Printf("Found %d tags for repository %s", len(tags), repositoryUrl)
	return tags, nil
}
