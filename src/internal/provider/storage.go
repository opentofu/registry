package provider

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
providerDirectoryRegex is a regular expression that matches the directory structure of a provider file.
  - (?i) makes the match case-insensitive.
  - providers/ matches the literal string "providers/".
  - \w matches a single word character. This corresponds to the first letter of the namespace.
  - (?P<Namespace>[^/]+) captures a sequence of one or more characters that are not a slash. This corresponds to "oracle".
  - (?P<ProviderName>[^/]+) captures another sequence of non-slash characters. This corresponds to "oci".
  - \.json matches the literal string ".json".
*/
var providerDirectoryMatcher = re.MustCompile(`(?i)providers/\w/(?P<Namespace>[^/]+?)/(?P<ProviderName>[^/]+?)\.json`)

func relative_path(id Identifier) string {
	return filepath.Join(strings.ToLower(id.Namespace[0:1]), id.Namespace, id.ProviderName+".json")
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
		Log:    log.With(slog.String("type", "provider")),
	}
}

func (s Storage) Create(id Identifier) Provider {
	return Provider{
		Identifier: id,
		Log: s.Log.With(
			slog.Group("provider", slog.String("namespace", id.Namespace), slog.String("name", id.ProviderName)),
		),
		Repository: s.Github.Repository(id.Namespace, fmt.Sprintf("terraform-provider-%s", id.ProviderName)),
	}
}

func (s Storage) Load(id Identifier) (Provider, error) {
	prov := s.Create(id)
	err := s.FS.ReadJSONFrom(relative_path(id), &prov.Metadata)
	return prov, err
}
func (s Storage) Save(prov Provider) error {
	return s.FS.WriteJSONInto(relative_path(prov.Identifier), prov.Metadata)
}

type Identifiers []Identifier

func (s Storage) List() (Identifiers, error) {
	paths, err := s.FS.List()
	if err != nil {
		return nil, err
	}
	ids := make([]Identifier, 0, len(paths))
	for _, path := range paths {
		match := providerDirectoryMatcher.Match(path)
		if match == nil {
			s.Log.Debug("Failed to extract provider details from path, skipping", slog.String("path", path))
			continue
		}
		ids = append(ids, Identifier{
			Namespace:    match["Namespace"],
			ProviderName: match["ProviderName"],
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
