package v1api

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/provider"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ProviderGenerator(t *testing.T) {
	logger := slog.Default()
	ctx := context.Background()

	p := ProviderGenerator{
		Provider: provider.Provider{
			Namespace:    "zededa",
			ProviderName: "zedcloud",
			Github:       github.NewClient(ctx, logger, ""),
			Logger:       logger,
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
	assert.Equal(t, "v2.3.1-RC1", list.Versions[0].Version)

	v := provider.Version{
		Version: "v2.3.1-RC1",
	}

	d := ProviderVersionDetails{
		Arch: "arm",
		OS:   "mac",
	}
	// Testing if the generated file is lower case
	assert.Equal(t, "gen/v1/providers/zededa/zedcloud/v2.3.1-rc1/download/mac/arm", p.VersionDownloadPath(v, d))

	vt, err := p.VersionFromTag("v2.3.1-RC1")
	require.NoError(t, err)

	// URLs still should have an uppercase version, since we point directly to Github URLs
	baseURL := "https://github.com/zededa/terraform-provider-zedcloud/releases/download/v2.3.1-RC1/terraform-provider-zedcloud_"
	expectedURL := fmt.Sprintf("%s%s", baseURL, "2.3.1-RC1_SHA256SUMS")
	assert.Equal(t, expectedURL, vt.SHASumsURL)

	expectedSigURL := fmt.Sprintf("%s%s", baseURL, "2.3.1-RC1_SHA256SUMS.sig")
	assert.Equal(t, expectedSigURL, vt.SHASumsSignatureURL)
}
