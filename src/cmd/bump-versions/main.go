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
	slog.SetDefault(logger) // TODO REMOVE ME
	logger.Info("Starting version bump process for modules and providers")

	moduleDataDir := flag.String("module-data", "../modules", "Directory containing the module data")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")

	flag.Parse()

	ctx := context.Background()
	token, err := github.EnvAuthToken()
	if err != nil {
		slog.Error("Initialization Error", slog.Any("err", err))
		os.Exit(1)
	}
	ghClient := github.NewClient(ctx, logger, token)

	modules, err := module.ListModules(*moduleDataDir)
	if err != nil {
		logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	for _, m := range modules {
		slog.Info("Beginning version bump process for module", slog.String("module", m.Namespace+"/"+m.Name+"/"+m.TargetSystem))
		err = module.UpdateMetadataFile(m, *moduleDataDir)
		if err != nil {
			slog.Error("Failed to version bump module", slog.Any("err", err))
			os.Exit(1)
		}
	}

	providers, err := provider.ListProviders(*providerDataDir, logger, ghClient)
	if err != nil {
		logger.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}

	errChan := make(chan error, len(providers))

	for _, p := range providers {
		p := p
		go func() {
			errChan <- p.UpdateMetadataFile()
		}()
	}

	var errs []error
	for _ = range providers {
		err := <-errChan
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		logger.Error("Encountered errors while processing providers")
		for _, err := range errs {
			logger.Error(err.Error())
		}
	}

	logger.Info("Completed version bump process for modules and providers")
}
