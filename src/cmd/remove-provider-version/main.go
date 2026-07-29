package main
 
import (
	"flag"
	"log/slog"
	"os"
 
	"github.com/opentofu/registry-stable/internal/provider"
)
 
// remove-provider-version deletes a single version entry from a provider's
// metadata file. It is intentionally separate from verify-reindex-request:
// that command only checks whether a re-index request is *allowed*, this
// command is the one that actually performs the mutation, and it is only
// ever invoked after verification has already succeeded.
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
 
	namespace := flag.String("namespace", "", "Provider namespace")
	name := flag.String("name", "", "Provider name")
	version := flag.String("version", "", "The provider version to remove")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")
	flag.Parse()
 
	if *namespace == "" || *name == "" || *version == "" {
		logger.Error("--namespace, --name, and --version are required")
		os.Exit(1)
	}
 
	p := provider.Provider{
		Namespace:    *namespace,
		ProviderName: *name,
		Directory:    *providerDataDir,
		Logger:       logger,
	}
 
	if err := p.RemoveVersion(*version); err != nil {
		logger.Error("Failed to remove version", slog.Any("err", err))
		os.Exit(1)
	}
 
	logger.Info("Removed version from metadata", slog.String("namespace", *namespace), slog.String("name", *name), slog.String("version", *version))
}
 
