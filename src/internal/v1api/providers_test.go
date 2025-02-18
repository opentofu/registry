package v1api

import (
	"log/slog"
	"testing"

	"github.com/opentofu/registry-stable/internal/provider"

	"github.com/stretchr/testify/assert"
)

func Test_ProviderGenerator(t *testing.T) {
	logger := slog.Default()

	p := ProviderGenerator{
		Provider: provider.Provider{
			Namespace:    "spacename",
			ProviderName: "name",
		},
		Metadata: provider.Metadata{
			Versions: []provider.Version{
				{
					Version: "v2.3.1-RC1",
				},
			},
		},
		Destination: "gen",
		log:         logger,
	}

	list := p.VersionListing()
	assert.Equal(t, "v2.3.1-rc1", list.Versions[0].Version)

	v := provider.Version{
		Version: "v2.3.1-RC1",
	}

	d := ProviderVersionDetails{
		Arch: "arm",
		OS:   "mac",
	}
	assert.Equal(t, "gen/v1/providers/spacename/name/v2.3.1-rc1/download/mac/arm", p.VersionDownloadPath(v, d))
}
