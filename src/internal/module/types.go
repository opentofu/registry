package module

import "fmt"

type Version struct {
	Version string `json:"version"` // The version number of the provider. Correlates to a tag in the module repository
}

type MetadataFile struct {
	Versions []Version `json:"versions"`
}

type Module struct {
	Namespace    string // The module namespace
	Name         string // The module name
	TargetSystem string // The module target system
}

func (m Module) RepositoryURL() string {
	repoName := fmt.Sprintf("terraform-%s-%s", m.TargetSystem, m.Name)
	return fmt.Sprintf("https://github.com/%s/%s", m.Namespace, repoName)
}
