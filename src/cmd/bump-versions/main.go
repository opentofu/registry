package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/opentofu/registry-stable/internal/blacklist"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/module"
	"github.com/opentofu/registry-stable/internal/provider"
)

func main() {
	started := time.Now()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Starting version bump process for modules and providers")

	moduleDataDir := flag.String("module-data", "../modules", "Directory containing the module data")
	moduleNamespace := flag.String("module-namespace", "", "Which module namespace to limit the command to")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")
	providerNamespace := flag.String("provider-namespace", "", "Which provider namespace to limit the command to")
	targetDuration := flag.Duration("target-duration", time.Minute*0, "Used to limit how much of a backfill this command can perform") // Backfill complete, disabled

	flag.Parse()

	ctx := context.Background()
	token, err := github.EnvAuthToken()
	if err != nil {
		logger.Error("Initialization Error", slog.Any("err", err))
		os.Exit(1)
	}
	ghClient := github.NewClient(ctx, logger, token)

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

	err = modules.Parallel(200, func(m module.Module) error {
		return m.UpdateMetadataFile()
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
	err = providers.Parallel(200, func(p provider.Provider) error {
		return p.UpdateMetadataFile()
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Completed version bump process for modules and providers")

	deadline := started.Add(*targetDuration)
	remainingTime := time.Until(deadline)

	ctx, cancel := context.WithTimeout(ctx, remainingTime)
	defer cancel()

	if ctx.Err() == context.DeadlineExceeded {
		logger.Info("Skipping backfill process, deadline exceeded")
		return
	}

	logger.Info("Beginning backfill process for providers", slog.Any("time_allocated", remainingTime))

	// Setup a new github client with the limited context
	ghClient = github.NewClient(ctx, logger, token)

	// Re-list providers with new github client (not ideal)
	providers, err = provider.ListProviders(*providerDataDir, *providerNamespace, logger, ghClient, bl)
	if err != nil {
		logger.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}
	err = providers.Parallel(20, func(p provider.Provider) error {
		if err := ctx.Err(); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				// Outta-time
				return nil
			}
			return err
		}
		err := p.BackfillVersionData(ctx)
		if errors.Is(err, context.DeadlineExceeded) {
			p.Logger.Info("Partial completion of backfill due to time limitation")
			// Outta-time
			return nil
		}
		return err
	})

	if err != nil {
		logger.Error("Failed to backfill providers", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Completed backfill process for providers")
}
