package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/module"
	"github.com/opentofu/registry-stable/internal/provider"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Starting version bump process for modules and providers")

	moduleDataDir := flag.String("module-data", "../modules", "Directory containing the module data")
	//moduleNamespace := flag.String("module-namespace", "", "Which module namespace to limit the command to")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")
	providerNamespace := flag.String("provider-namespace", "", "Which provider namespace to limit the command to")

	flag.Parse()

	ctx := context.Background()
	token, err := github.EnvAuthToken()
	if err != nil {
		logger.Error("Initialization Error", slog.Any("err", err))
		os.Exit(1)
	}
	ghClient := github.NewClient(ctx, logger, token)

	modStorage := module.NewStorage(*moduleDataDir, logger, ghClient)
	modIds, err := modStorage.List()
	if err != nil {
		logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	errs := ForEachModuleInParallel(modIds, func(id module.Identifier) error {
		m, err := modStorage.Load(id)
		if err != nil {
			return err
		}
		err = m.UpdateMetadata()
		if err != nil {
			return err
		}
		return modStorage.Save(m)
	})

	if len(errs) != 0 {
		logger.Error("Errors occured while processing modules")
		for _, err := range errs {
			logger.Error(err.Error())
		}
		os.Exit(1)
	}

	providers, err := provider.ListProviders(*providerDataDir, *providerNamespace, logger, ghClient)
	if err != nil {
		logger.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}
	err = providers.Parallel(200, func(p provider.Provider) error {
		return p.UpdateMetadataFile()
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Completed version bump process for modules and providers")
}

func ForEachModuleInParallel(modIds []module.Identifier, fn func(module.Identifier) error) []error {
	errChan := make(chan error, len(modIds))
	for _, id := range modIds {
		id := id
		go func() {
			errChan <- fn(id)
		}()
	}

	errs := make([]error, 0)
	for range modIds {
		err := <-errChan
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
