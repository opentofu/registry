package module

import (
	"os"
	"path/filepath"
	"regexp"

	"registry-stable/internal/module"
)

func ListModules(moduleDataDir string) ([]module.Module, error) {
	// walk the module directory recursively and find all json files
	// for each json file, parse it into a module.Module struct
	// return a slice of module.Module structs

	var results []module.Module
	err := filepath.Walk(moduleDataDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".json" {
			// extract the module details from the file path. It should be modules/<firstLetterOfNamespace>/<namespace>/<name>/<system>.json
			regex := regexp.MustCompile(`modules\/([a-z])\/([a-z0-9-]+)\/([a-z0-9-]+)\/([a-z0-9-]+).json`)
			matches := regex.FindStringSubmatch(path)
			if len(matches) == 5 {
				results = append(results, module.Module{
					Namespace:    matches[2],
					Name:         matches[3],
					TargetSystem: matches[4],
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}
