package module

type Version struct {
	Version string `json:"version"` // The version number of the provider. Correlates to a tag in the module repository
}

type MetadataFile struct {
	Versions []Version `json:"versions"`
}

type Module struct {
	Namespace string // The module namespace
	Name      string // The module name
	System    string // The module system
}
