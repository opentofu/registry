package main

import (
	"flag"
	"log/slog"
	"os"

	"registry-stable/cmd/common"
	"registry-stable/internal/module"
	"registry-stable/internal/provider"
	"registry-stable/internal/v1api"
)

func main() {
	destinationDir := flag.String("destination", "../generated", "Directory to write the generated responses to")

	cli := common.Parse()

	cli.Logger.Info("Generating v1 API responses")

	err := v1api.WriteWellKnownFile(*destinationDir)
	if err != nil {
		cli.Logger.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	err = cli.Modules(func(m module.Module) error {
		g, err := v1api.NewModuleGenerator(m, *destinationDir)
		if err != nil {
			return err
		}

		return g.Generate()
	})
	if err != nil {
		cli.Logger.Error("Failed to process modules", slog.Any("err", err))
		os.Exit(1)
	}

	err = cli.Providers(func(p provider.Provider) error {
		g, err := v1api.NewProviderGenerator(p, *destinationDir)
		if err != nil {
			return err
		}
		return g.Generate()
	})
	if err != nil {
		cli.Logger.Error("Failed to process modules", slog.Any("err", err))
		os.Exit(1)
	}

	cli.Logger.Info("Completed generating v1 API responses")
}
