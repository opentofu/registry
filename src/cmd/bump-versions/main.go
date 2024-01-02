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
	moduleNamespace := flag.String("module-namespace", "", "Which module namespace to limit the command to")
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
	modules, err := modStorage.List()
	if err != nil {
		logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	errs := modules.ParallelForEach(func(id module.Identifier) error {
		if *moduleNamespace != "" && *moduleNamespace != id.Namespace {
			return nil
		}

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

	provStorage := provider.NewStorage(*providerDataDir, logger, ghClient)
	providers, err := provStorage.List()
	errs = providers.ParallelForEach(func(id provider.Identifier) error {
		if *providerNamespace != "" && *providerNamespace != id.Namespace {
			return nil
		}

		p, err := provStorage.Load(id)
		if err != nil {
			return err
		}
		err = p.UpdateMetadata()
		if err != nil {
			return err
		}
		return provStorage.Save(p)
	})
	if len(errs) != 0 {
		logger.Error("Errors occured while processing providers")
		for _, err := range errs {
			logger.Error(err.Error())
		}
		os.Exit(1)
	}

	logger.Info("Completed version bump process for modules and providers")
}
