package module

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/parallel"
)

/*
moduleDirectoryRegex is a regular expression that matches the directory structure of a module file.
  - (?i) makes the match case-insensitive.
  - modules/ matches the literal string "modules/".
  - \w matches a single word character. This corresponds to the first letter of the namespace.
  - (?P<Namespace>[^/]+) captures a sequence of one or more characters that are not a slash. This corresponds to "terraform-aws-modules".
  - (?P<Name>[^/]+) captures another sequence of non-slash characters. This corresponds to "lambda".
  - (?P<TargetSystem>[^/]+) captures the third sequence of non-slash characters. This corresponds to "aws".
  - \.json matches the literal string ".json".
*/
var moduleDirectoryRegex = regexp.MustCompile(`(?i)modules/\w/(?P<Namespace>[^/]+?)/(?P<Name>[^/]+?)/(?P<TargetSystem>[^/]+?)\.json`)

func extractModuleDetailsFromPath(path string) *Module {
	matches := moduleDirectoryRegex.FindStringSubmatch(path)
	if len(matches) != 4 {
		return nil
	}

	m := Module{
		Namespace:    matches[moduleDirectoryRegex.SubexpIndex("Namespace")],
		Name:         matches[moduleDirectoryRegex.SubexpIndex("Name")],
		TargetSystem: matches[moduleDirectoryRegex.SubexpIndex("TargetSystem")],
	}

	return &m
}

type ModuleList []Module

func ListModules(moduleDataDir string, logger *slog.Logger, ghClient github.Client) (ModuleList, error) {
	// walk the module directory recursively and find all json files
	// for each json file, parse it into a Module struct
	// return a slice of Module structs

	var results []Module
	err := filepath.Walk(moduleDataDir, func(path string, info os.FileInfo, err error) error {
		m := extractModuleDetailsFromPath(path)
		if m != nil {
			m.Directory = moduleDataDir
			m.Logger = logger.With(
				slog.String("type", "module"),
				slog.Group("module", slog.String("namespace", m.Namespace), slog.String("name", m.Name), slog.String("targetsystem", m.TargetSystem)),
			)
			m.Github = ghClient.WithLogger(m.Logger)
			results = append(results, *m)
		} else {
			logger.Debug("Failed to extract module details from path, skipping", slog.String("path", path))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

type Action func(m Module) error

func (modules ModuleList) Parallel(maxConcurrency int, action Action) error {
	actions := make([]parallel.Action, len(modules))
	for i, m := range modules {
		m := m
		actions[i] = func() error {
			err := action(m)
			if err != nil {
				m.Logger.Error(err.Error())
				return err
			}
			return nil
		}
	}

	errs := parallel.ForEach(actions, maxConcurrency)
	if len(errs) != 0 {
		return fmt.Errorf("Encountered %d errors processing %d modules", len(errs), len(modules))
	}
	return nil
}
