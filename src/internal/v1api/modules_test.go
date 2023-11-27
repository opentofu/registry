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
		MetadataFile: module.MetadataFile{},
		Destination:  "gen",
		log:          logger,
	}
	v := module.Version{
		Version: "v1.0.1",
	}
	assert.Equal(t, "gen/v1/modules/spacename/name/target/versions", m.VersionListingPath())
	assert.Equal(t, "gen/v1/modules/spacename/name/target/1.0.1/download", m.VersionDownloadPath(v))
}
