package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
