package module

import (
	"testing"
)

func TestValidate(t *testing.T) {
	type TestCase struct {
		name    string
		input   Metadata
		wantErr bool
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
			wantErr: true,
		},
		{
			name:    "empty-versions-list",
			input:   Metadata{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.input); (err != nil) != tt.wantErr {
				t.Fatalf("Validate(%v) unexpected error = %v", tt.input, err)
			}
		})
	}
}
