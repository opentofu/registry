package v1api

import (
	"registry-stable/internal/module"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ModuleGenerator(t *testing.T) {
	m := ModuleGenerator{
		module.Module{
			Namespace:    "spacename",
			Name:         "name",
			TargetSystem: "target",
		},
		module.MetadataFile{},
		"gen",
	}
	v := module.Version{
		Version: "v1.0.1",
	}
	assert.Equal(t, "gen/v1/modules/spacename/name/target/versions", m.VersionListingPath())
	assert.Equal(t, "gen/v1/modules/spacename/name/target/1.0.1/download", m.VersionDownloadPath(v))
}
