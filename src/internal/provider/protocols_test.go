package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func stringptr(s string) *string {
	return &s
}

func TestParseManifestContents(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		want        *Manifest
		expectedErr *string
	}{
		{
			name:        "Empty manifest",
			input:       []byte(``),
			expectedErr: stringptr("failed to parse manifest contents"),
		},
		{
			name:        "Invalid JSON",
			input:       []byte(`{"version":1,`),
			expectedErr: stringptr("failed to parse manifest contents"),
		},
		{
			name:  "Valid manifest",
			input: []byte(`{"version":1,"metadata":{"protocol_versions":["5.0"]}}`),
			want:  &Manifest{Metadata: ManifestMetadata{ProtocolVersions: []string{"5.0"}}},
		},
		{
			name:  "Valid manifest - multiple protocol versions",
			input: []byte(`{"version":1,"metadata":{"protocol_versions":["5.0", "5.1"]}}`),
			want:  &Manifest{Metadata: ManifestMetadata{ProtocolVersions: []string{"5.0", "5.1"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseManifestContents(tt.input)
			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, *tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
