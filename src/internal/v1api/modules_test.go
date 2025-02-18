package v1api

import (
	"log/slog"
	"testing"

	"github.com/opentofu/registry-stable/internal/module"

	"github.com/stretchr/testify/assert"
)

func Test_ModuleGenerator(t *testing.T) {
	logger := slog.Default()

	m := ModuleGenerator{
		Module: module.Module{
			Namespace:    "spacename",
			Name:         "name",
			TargetSystem: "target",
		},
		Metadata:    module.Metadata{},
		Destination: "gen",
		log:         logger,
	}
	v := module.Version{
		Version: "v1.0.1",
	}

	vl := module.Version{
		Version: "v2.3.1-RC1",
	}
	assert.Equal(t, "gen/v1/modules/spacename/name/target/versions", m.VersionListingPath())
	assert.Equal(t, "gen/v1/modules/spacename/name/target/1.0.1/download", m.VersionDownloadPath(v))
	assert.Equal(t, "gen/v1/modules/spacename/name/target/2.3.1-rc1/download", m.VersionDownloadPath(vl))
}
