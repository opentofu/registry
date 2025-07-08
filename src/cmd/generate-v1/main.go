package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/opentofu/registry-stable/internal/blacklist"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/module"
	"github.com/opentofu/registry-stable/internal/provider"
	"github.com/opentofu/registry-stable/internal/v1api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Generating v1 API responses")

	moduleDataDir := flag.String("module-data", "../modules", "Directory containing the module data")
	moduleNamespace := flag.String("module-namespace", "", "Which module namespace to limit the command to")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")
	providerNamespace := flag.String("provider-namespace", "", "Which provider namespace to limit the command to")
	keyDataDir := flag.String("key-data", "../keys", "Directory containing the gpg keys")

	destinationDir := flag.String("destination", "../generated", "Directory to write the generated responses to")

	// Will panic if used, it should not be.
	// In the future we probably want to change module.Module/module.Meta -> module.Identifer/module.Module
	ghClient := github.Client{}

	flag.Parse()

	err := v1api.WriteWellKnownFile(*destinationDir)
	if err != nil {
		logger.Error("Failed to create well known file", slog.Any("err", err))
		os.Exit(1)
	}
	bl, err := blacklist.Load()
	if err != nil {
		logger.Error("Failed to load blacklist", slog.Any("err", err))
		os.Exit(1)
	} else {
		logger.Info("Loaded blacklist successfully")
	}

	modules, err := module.ListModules(*moduleDataDir, *moduleNamespace, logger, ghClient, bl)
	if err != nil {
		logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}
	err = modules.Parallel(20, func(m module.Module) error {
		g, err := v1api.NewModuleGenerator(m, *destinationDir)
		if err != nil {
			return err
		}
		return g.Generate()
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	providers, err := provider.ListProviders(*providerDataDir, *providerNamespace, logger, ghClient, bl)
	if err != nil {
		logger.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}
	err = providers.Parallel(20, func(p provider.Provider) error {
		g, err := v1api.NewProviderGenerator(p, *destinationDir, *keyDataDir)
		if err != nil {
			return err
		}
		return g.Generate()
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if *providerNamespace == "" || *providerNamespace == "hashicorp" {
		err = v1api.ArchivedOverrides(*destinationDir, logger)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	logger.Info("Completed generating v1 API responses")
}
