package github

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// TODO logs
// TODO Dependency Inversion for tests?

func GetTags(repositoryUrl string) ([]string, error) {
	log.Printf("Getting tags for repository %s", repositoryUrl)

	var buf bytes.Buffer
	var bufErr bytes.Buffer
	cmd := exec.Command("git", "ls-remote", "--tags", "--refs", repositoryUrl)
	cmd.Stdout = &buf
	cmd.Stderr = &bufErr

	if err := cmd.Run(); err != nil {
		log.Printf("Could not get tags for repository %s: %s", repositoryUrl, bufErr.String())
		return nil, err
	}

	tags := make([]string, 0)
	for _, line := range strings.Split(buf.String(), "\n") {
		if !strings.Contains(line, "refs/tags/") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("could not parse tags: tags are in wrong format")
		}

		ref := fields[1]
		if !strings.HasPrefix(ref, "refs/tags/") {
			continue
		}

		tag := strings.TrimPrefix(ref, "refs/tags/")
		tags = append(tags, tag)
	}

	log.Printf("Found %d tags for repository %s", len(tags), repositoryUrl)

	return tags, nil
}
