package github

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseTagsFromStdout(t *testing.T) {
	cases := map[string]struct {
		input       []string
		expected    []string
		expectedErr error
	}{
		"Empty input": {
			input:    []string{""},
			expected: []string{},
		},
		// Successful
		"Simple Tag": {
			input:    []string{"314159265358979     refs/tags/v0.0.1"},
			expected: []string{"v0.0.1"},
		},
		"Multiple Tags": {
			input:    []string{"314159265358979     refs/tags/v0.0.1", "314159265358979     refs/tags/v0.1.1", "314159265358979     refs/tags/v1.0.1"},
			expected: []string{"v0.0.1", "v0.1.1", "v1.0.1"},
		},
		// Invalid entries (ignored)
		"No Tags": {
			input:    []string{},
			expected: []string{},
		},
		"Empty Tags": {
			input:    []string{""},
			expected: []string{},
		},
		"Multiple Tags w/ Invalid": {
			input:    []string{"314159265358979     HEAD", "314159265358979     refs/tags/v0.1.1", "314159265358979     refs/tags/v1.0.1"},
			expected: []string{"v0.1.1", "v1.0.1"},
		},
		// Error cases
		"Missing Field": {
			input:       []string{"borkborkborkrefs/tags/"},
			expectedErr: errors.New("invalid format for tag 'borkborkborkrefs/tags/', expected two fields"),
		},
		"Extra Field": {
			input:       []string{"bork bork refs/tags/"},
			expectedErr: errors.New("invalid format for tag 'bork bork refs/tags/', expected two fields"),
		},
		"Missing tags/refs": {
			input:       []string{"314159265358979refs/tags/   v0.0.1"},
			expectedErr: errors.New("invalid format for tag '314159265358979refs/tags/   v0.0.1', expected 'refs/tags/' prefix"),
		},
		"Missing version": {
			input:       []string{"314159265358979 refs/tags/"},
			expectedErr: errors.New("invalid format for tag '314159265358979 refs/tags/', no version provided"),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := parseTagsFromStdout(tc.input)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, out)

		})
	}
}
