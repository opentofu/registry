package github

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parseTagsFromStdout(t *testing.T) {
	cases := map[string]struct {
		input       []string
		expected    []string
		expectedErr string
	}{
		"Simple Tag": {
			input:    []string{"314159265358979     refs/tags/v0.0.1"},
			expected: []string{"v0.0.1"},
		},
		"Multiple Tags": {
			input:    []string{"314159265358979     refs/tags/v0.0.1", "314159265358979     refs/tags/v0.1.1", "314159265358979     refs/tags/v1.0.1"},
			expected: []string{"v0.0.1", "v0.1.1", "v1.0.1"},
		},
		"Multiple Tags w/ Invalid": {
			input:    []string{"314159265358979     HEAD", "314159265358979     refs/tags/v0.1.1", "314159265358979     refs/tags/v1.0.1"},
			expected: []string{"v0.1.1", "v1.0.1"},
		},
		"Missing Field": {
			input:       []string{"borkborkborkrefs/tags/"},
			expectedErr: "Invalid format for tag 'borkborkborkrefs/tags/', expected two fields",
		},
		"Extra Field": {
			input:       []string{"bork bork refs/tags/"},
			expectedErr: "Invalid format for tag 'bork bork refs/tags/', expected two fields",
		},
		"Missing tags/refs": {
			input:       []string{"314159265358979refs/tags/   v0.0.1"},
			expectedErr: "Invalid format for tag '314159265358979refs/tags/   v0.0.1', expected 'refs/tags/' prefix",
		},
		"Missing version": {
			input:       []string{"314159265358979 refs/tags/"},
			expectedErr: "Invalid format for tag '314159265358979 refs/tags/', no version provided",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := parseTagsFromStdout(tc.input)

			if tc.expectedErr == "" {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error %s, got <nil>", tc.expectedErr)
				} else if err.Error() != tc.expectedErr {
					t.Errorf("Expected error %s, got %v", tc.expectedErr, err)
				}
			}

			if diff := cmp.Diff(out, tc.expected); diff != "" {
				t.Errorf("Expected %v, got %v: %s", tc.expected, out, diff)
			}

		})
	}
}
