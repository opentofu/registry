package main

import (
	"context"
	"log/slog"
	"os"
	"registry-stable/cmd/common"
	"registry-stable/internal/github"
	"registry-stable/internal/module"
	"registry-stable/internal/provider"
)

func main() {
	cli := common.Parse()

	cli.Logger.Info("Starting version bump process for modules and providers")

	ctx := context.Background()
	token, err := github.EnvAuthToken()
	if err != nil {
		slog.Error("Initialization Error", slog.Any("err", err))
		os.Exit(1)
	}
	cli.Github = github.NewClient(ctx, cli.Logger, token)

	err = cli.Modules(func(m module.Module) error { return m.UpdateMetadataFile() })
	if err != nil {
		cli.Logger.Error("Failed to process modules", slog.Any("err", err))
		os.Exit(1)
	}

	err = cli.Providers(func(p provider.Provider) error { return p.UpdateMetadataFile() })
	if err != nil {
		cli.Logger.Error("Failed to process providers", slog.Any("err", err))
		os.Exit(1)
	}

	cli.Logger.Info("Completed version bump process for modules and providers")
}
