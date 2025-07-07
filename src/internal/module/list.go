package module

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/opentofu/registry-stable/internal/blacklist"
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

// ListModules walks the module metadata directory provided and returns a list of modules.
func ListModules(moduleDataDir string, moduleNamespace string, logger *slog.Logger, ghClient github.Client, blacklistInstance *blacklist.Blacklist) (List, error) {
	var results []Module
	moduleNamespace = strings.ToLower(moduleNamespace)
	err := filepath.Walk(moduleDataDir, func(path string, info os.FileInfo, err error) error {
		m := extractModuleDetailsFromPath(path)
		if m == nil {
			logger.Debug("Failed to extract module details from path, skipping", slog.String("path", path))
			return nil
		}

		if moduleNamespace != "" && !strings.HasPrefix(strings.ToLower(m.Namespace), moduleNamespace) {
			logger.Debug("Filtered out module", slog.String("path", path))
			return nil
		}

		// enrich the module with additional information
		m.Directory = moduleDataDir
		m.Logger = logger.With(
			slog.String("type", "module"),
			slog.Group("module", slog.String("namespace", m.Namespace), slog.String("name", m.Name), slog.String("targetsystem", m.TargetSystem)),
		)
		m.Github = ghClient.WithLogger(m.Logger)
		m.Blacklist = blacklistInstance
		results = append(results, *m)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// List is a list of modules
type List []Module

// Action is a function that takes a module and returns an error, to be used by List.Parallel
type Action func(m Module) error

// Parallel runs the given action on each module in the list in parallel by wrapping parallel.ForEach
func (l List) Parallel(maxConcurrency int, action Action) error {
	actions := make([]parallel.Action, len(l))
	for i, m := range l {
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
		return fmt.Errorf("encountered %d errors processing %d modules", len(errs), len(l))
	}
	return nil
}
