package module

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/opentofu/registry-stable/internal/base"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/parallel"
	"github.com/opentofu/registry-stable/internal/re"
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
var moduleDirectoryMatcher = re.MustCompile(`(?i)modules/\w/(?P<Namespace>[^/]+?)/(?P<Name>[^/]+?)/(?P<TargetSystem>[^/]+?)\.json`)

func relative_path(id Identifier) string {
	return filepath.Join(strings.ToLower(id.Namespace[0:1]), id.Namespace, id.Name, id.TargetSystem+".json")
}

type Storage struct {
	FS     base.FileSystem
	Github github.Client
	Log    *slog.Logger
}

func NewStorage(directory string, log *slog.Logger, client github.Client) Storage {
	return Storage{
		FS:     base.FileSystem{Directory: directory},
		Github: client,
		Log:    log.With(slog.String("type", "module")),
	}
}

func (s Storage) Create(id Identifier) Module {
	log := s.Log.With(slog.Group("module",
		slog.String("namespace", id.Namespace),
		slog.String("name", id.Name),
		slog.String("targetsystem", id.TargetSystem),
	))
	return Module{
		Identifier: id,
		Log:        log,
		Repository: s.Github.Repository(id.Namespace, fmt.Sprintf("terraform-%s-%s", id.TargetSystem, id.Name), log),
	}
}

func (s Storage) Load(id Identifier) (Module, error) {
	mod := s.Create(id)
	err := s.FS.ReadJSONFrom(relative_path(id), &mod.Metadata)
	return mod, err
}
func (s Storage) Save(mod Module) error {
	return s.FS.WriteJSONInto(relative_path(mod.Identifier), mod.Metadata)
}

type Identifiers []Identifier

func (s Storage) List() (Identifiers, error) {
	paths, err := s.FS.List()
	if err != nil {
		return nil, err
	}
	ids := make(Identifiers, 0, len(paths))
	for _, path := range paths {
		match := moduleDirectoryMatcher.Match(path)
		if match == nil {
			s.Log.Debug("Failed to extract module details from path, skipping", slog.String("path", path))
			continue
		}
		ids = append(ids, Identifier{
			Namespace:    match["Namespace"],
			Name:         match["Name"],
			TargetSystem: match["TargetSystem"],
		})
	}

	return ids, nil
}

func (l Identifiers) ParallelForEach(action func(Identifier) error) []error {
	eg := make(parallel.ErrorGroup, 0)
	for _, t := range l {
		t := t
		eg = append(eg, func() error {
			return action(t)
		})
	}
	return eg.Errors()
}
