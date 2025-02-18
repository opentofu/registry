package v1api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/provider"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/v2.3.1-RC1/terraform-provider-name_2.3.1-RC1_SHA256SUMS", func(w http.ResponseWriter, r *http.Request) {
		checksums := []byte(`
123 asd.zip
121 terraform-provider-name_v2.3.1-RC1_darwin_386.zip
`)
		_, err := w.Write(checksums)
		if err != nil {
			t.Fatalf("Couldn't write to testing response of /: %v", err)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("/"))
		if err != nil {
			t.Fatalf("Couldn't write to testing response of /: %v", err)
		}
	})

	srv := httptest.NewServer(mux)

	t.Cleanup(func() {
		srv.Close()
	})
	return srv
}

func Test_ProviderGenerator(t *testing.T) {
	logger := slog.Default()
	srv := newTestServer(t)
	ctx := context.Background()

	p := ProviderGenerator{
		Provider: provider.Provider{
			Namespace:    "spacename",
			ProviderName: "name",
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
	assert.Equal(t, "gen/v1/providers/spacename/name/v2.3.1-rc1/download/mac/arm", p.VersionDownloadPath(v, d))

	args := provider.VersionFromTagArgs{
		URLPrefix: srv.URL,
		Release:   "v2.3.1-RC1",
	}
	vt, err := p.VersionFromTag(args)
	require.NoError(t, err)

	expectedURL := fmt.Sprintf("%s/%s/%s", srv.URL, "v2.3.1-RC1", "terraform-provider-name_2.3.1-RC1_SHA256SUMS")
	assert.Equal(t, expectedURL, vt.SHASumsURL)

	expectedSigURL := fmt.Sprintf("%s/%s/%s", srv.URL, "v2.3.1-RC1", "terraform-provider-name_2.3.1-RC1_SHA256SUMS.sig")
	assert.Equal(t, expectedSigURL, vt.SHASumsSignatureURL)
}
