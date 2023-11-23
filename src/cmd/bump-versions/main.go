package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"registry-stable/internal/github"
	"registry-stable/internal/module"
	"registry-stable/internal/provider"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Starting version bump process for modules and providers")

	moduleDataDir := flag.String("module-data", "../modules", "Directory containing the module data")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")

	flag.Parse()

	ctx := context.Background()
	token, err := github.EnvAuthToken()
	if err != nil {
		logger.Error("Initialization Error", slog.Any("err", err))
		os.Exit(1)
	}
	ghClient := github.NewClient(ctx, logger, token)

	modules, err := module.ListModules(*moduleDataDir, logger, ghClient)
	if err != nil {
		logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	err = modules.Parallel(20, func(m module.Module) error {
		return m.UpdateMetadataFile()
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	providers, err := provider.ListProviders(*providerDataDir, logger, ghClient)
	if err != nil {
		logger.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}
	err = providers.Parallel(20, func(p provider.Provider) error {
		return p.UpdateMetadataFile()
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Completed version bump process for modules and providers")
}
