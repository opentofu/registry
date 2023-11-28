package github

import (
	"log/slog"
	"testing"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestExtractTag(t *testing.T) {

	stringPtr := func(s string) *string {
		return &s
	}

	tests := []struct {
		name         string
		guid         string
		expectedTag  *string
		expectingErr bool
	}{
		{
			name:        "Valid tag with version v5.21.0",
			guid:        "tag:github.com,2008:Repository/72815297/v5.21.0",
			expectedTag: stringPtr("v5.21.0"),
		},
		{
			name:        "Valid tag with version v5.12.0",
			guid:        "tag:github.com,2008:Repository/72815297/v5.12.0",
			expectedTag: stringPtr("v5.12.0"),
		},
		{
			name:        "Invalid tag with no version",
			guid:        "tag:github.com,2008:Repository/72815297/",
			expectedTag: nil,
		},
		{
			name:        "empty tag",
			guid:        "",
			expectedTag: nil,
		},
	}

	client := Client{
		log: slog.Default(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// construct a fake gofeed item to be consumed by extractTag
			// we only care about the GUID field
			item := &gofeed.Item{
				GUID: tt.guid,
			}
			tag := client.extractTag(item)
			assert.Equal(t, tt.expectedTag, tag)
		})
	}
}
