package module

func UpdateMetadataFile(m Module, moduleDataDir string) error {
	// No need to check whether there's a new version or not
	// Simply recreate the metadata file from scratch. If there's no change, it will remain as-is
	return CreateMetadataFile(m, moduleDataDir)
}
