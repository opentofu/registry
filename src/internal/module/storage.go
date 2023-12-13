package module

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/opentofu/registry-stable/internal/base"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/re"
)

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
	return Module{
		Identifier: id,
		Log:        s.Log,
		Repository: s.Github.Repository(id.Namespace, fmt.Sprintf("terraform-%s-%s", id.TargetSystem, id.Name)),
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

func (s Storage) List() ([]Identifier, error) {
	paths, err := s.FS.List()
	if err != nil {
		return nil, err
	}
	ids := make([]Identifier, 0, len(paths))
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
