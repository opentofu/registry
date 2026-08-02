package github

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func installFakeGit(t *testing.T, body string) {
	t.Helper()

	dir := t.TempDir()
	script := `#!/bin/sh
if [ "$GIT_TERMINAL_PROMPT" != "0" ]; then
	echo "unexpected GIT_TERMINAL_PROMPT=$GIT_TERMINAL_PROMPT" >&2
	exit 1
fi
` + body
	require.NoError(t, os.WriteFile(filepath.Join(dir, "git"), []byte(script), 0755))
	t.Setenv("PATH", dir)
}

func newTestClient() Client {
	return Client{
		log: slog.New(slog.NewTextHandler(io.Discard, nil)),
		cliThrottle: func() ThrottleToken {
			return func() {}
		},
	}
}

func Test_parseTagsFromStdout(t *testing.T) {
	cases := map[string]struct {
		input       []string
		expected    []Tag
		expectedErr string
	}{
		"Empty input": {
			input:    []string{""},
			expected: []Tag{},
		},
		// Successful
		"Simple Tag": {
			input: []string{"3141592653589793     refs/tags/v0.0.1"},
			expected: []Tag{
				{Commit: "3141592653589793", Ref: "v0.0.1"},
			},
		},
		"Multiple Tags": {
			input: []string{"3141592653589793     refs/tags/v0.0.1", "3141592653589793     refs/tags/v0.1.1", "3141592653589793     refs/tags/v1.0.1"},
			expected: []Tag{
				{Commit: "3141592653589793", Ref: "v0.0.1"},
				{Commit: "3141592653589793", Ref: "v0.1.1"},
				{Commit: "3141592653589793", Ref: "v1.0.1"},
			},
		},
		// Invalid entries (ignored)
		"No Tags": {
			input:    []string{},
			expected: []Tag{},
		},
		"Empty Tags": {
			input:    []string{""},
			expected: []Tag{},
		},
		"Multiple Tags w/ Invalid": {
			input: []string{"3141592653589793     HEAD", "3141592653589793     refs/tags/v0.1.1", "3141592653589793     refs/tags/v1.0.1"},
			expected: []Tag{
				{Commit: "3141592653589793", Ref: "v0.1.1"},
				{Commit: "3141592653589793", Ref: "v1.0.1"},
			},
		},
		// Error cases
		"Missing Field": {
			input:       []string{"borkborkborkrefs/tags/"},
			expectedErr: "invalid format for tag \"borkborkborkrefs/tags/\", expected two fields",
		},
		"Extra Field": {
			input:       []string{"deadbeef deadbeef refs/tags/"},
			expectedErr: "invalid format for tag \"deadbeef deadbeef refs/tags/\", expected two fields",
		},
		"Bad commit": {
			input:       []string{"borkbork refs/tags/v0.0.1"},
			expectedErr: "invalid format for commit \"borkbork refs/tags/v0.0.1\": encoding/hex: invalid byte: U+006F 'o'",
		},
		"Missing tags/refs": {
			input:       []string{"3141592653589793   v0.0.1refs/tags/"},
			expectedErr: "invalid format for tag \"3141592653589793   v0.0.1refs/tags/\", expected \"refs/tags/\" prefix",
		},
		"Missing version": {
			input:       []string{"3141592653589793 refs/tags/"},
			expectedErr: "invalid format for tag \"3141592653589793 refs/tags/\", no version provided",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := parseTagsFromStdout(tc.input)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.expected, out)

		})
	}
}

func TestClientGetTagsDisablesGitTerminalPrompt(t *testing.T) {
	t.Setenv("GET_TAGS_ENV_SENTINEL", "inherited")
	installFakeGit(t, `if [ "$GET_TAGS_ENV_SENTINEL" != "inherited" ]; then
	echo "parent environment was not inherited" >&2
	exit 1
fi
printf '3141592653589793 refs/tags/v0.0.1\n'
`)

	tags, err := newTestClient().GetTags("https://example.com/example/repository.git")

	require.NoError(t, err)
	assert.Equal(t, []Tag{{Commit: "3141592653589793", Ref: "v0.0.1"}}, tags)
}

func TestClientGetTagsOverridesGitTerminalPrompt(t *testing.T) {
	t.Setenv("GIT_TERMINAL_PROMPT", "1")
	installFakeGit(t, "printf '3141592653589793 refs/tags/v0.0.1\\n'\n")

	tags, err := newTestClient().GetTags("https://example.com/example/repository.git")

	require.NoError(t, err)
	assert.Equal(t, []Tag{{Commit: "3141592653589793", Ref: "v0.0.1"}}, tags)
}

func TestClientGetTagsReturnsGitError(t *testing.T) {
	installFakeGit(t, `echo "fake git failure" >&2
exit 23
`)

	tags, err := newTestClient().GetTags("https://example.com/example/missing.git")

	assert.Nil(t, tags)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "could not get tags for https://example.com/example/missing.git")
	assert.Contains(t, err.Error(), "fake git failure")
}
