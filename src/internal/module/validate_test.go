package module

import (
	"testing"
)

func TestValidate(t *testing.T) {
	type TestCase struct {
		name       string
		input      Metadata
		wantErrStr string
	}

	tests := []TestCase{
		{
			name: "valid",
			input: Metadata{
				Versions: []Version{{"0.0.2"}, {"0.0.1"}},
			},
		},
		{
			name: "invalid-version",
			input: Metadata{
				Versions: []Version{{"0.0.2"}, {"foo"}},
			},
			wantErrStr: "found semver-incompatible version: foo\n",
		},
		{
			name:       "empty-versions-list",
			input:      Metadata{},
			wantErrStr: "found empty list of versions\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			switch tt.wantErrStr != "" {
			case true:
				if err == nil || tt.wantErrStr != err.Error() {
					t.Fatalf("unexpected error message, want = %s, got = %v", tt.wantErrStr, err)
				}
			default:
				if err != nil {
					t.Fatalf("unexpected error message: %v", err)
				}
			}
		})
	}
}
