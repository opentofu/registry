package main

import (
	"flag"
	"log/slog"
	"os"

	"registry-stable/internal/github"
	"registry-stable/internal/module"
	"registry-stable/internal/provider"
	"registry-stable/internal/v1api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Generating v1 API responses")

	moduleDataDir := flag.String("module-data", "../modules", "Directory containing the module data")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")
	destinationDir := flag.String("destination", "../generated", "Directory to write the generated responses to")

	// Will panic if used, it should not be.
	// In the future we probably want to change module.Module/module.Meta -> module.Identifer/module.Module
	ghClient := github.Client{}

	flag.Parse()

	err := v1api.WriteWellKnownFile(*destinationDir)
	if err != nil {
		logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	modules, err := module.ListModules(*moduleDataDir, logger, ghClient)
	if err != nil {
		logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	for _, m := range modules {
		g, err := v1api.NewModuleGenerator(m, *destinationDir)
		if err != nil {
			logger.Error("Failed to generate module version listing response", slog.Any("err", err))
			os.Exit(1)
		}

		err = g.Generate()
		if err != nil {
			logger.Error("Failed to generate module version listing response", slog.Any("err", err))
			os.Exit(1)
		}
	}

	providers, err := provider.ListProviders(*providerDataDir, logger, ghClient)
	if err != nil {
		slog.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}

	for _, p := range providers {
		g, err := v1api.NewProviderGenerator(p, *destinationDir)
		if err != nil {
			p.Logger.Error("Failed to generate provider version listing response", slog.Any("err", err))
			os.Exit(1)
		}

		err = g.Generate()
		if err != nil {
			p.Logger.Error("Failed to generate provider version listing response", slog.Any("err", err))
			os.Exit(1)
		}
	}

	slog.Info("Completed generating v1 API responses")
}
