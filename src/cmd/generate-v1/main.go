package main

import (
	"context"
	"log/slog"
	"os"

	"registry-stable/internal/files"
	"registry-stable/internal/repository-metadata-files/module"
	"registry-stable/internal/v1api"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.Info("Generating v1 API responses")

	ctx := context.Background()

	fs := os.DirFS(".")

	v1APIGenerator := v1api.Generator{
		ModuleDataDir:   "../modules",
		ProviderDataDir: "../providers",

		DestinationDir: "./generated",

		Filesystem: fs,
		FileWriter: &files.RealFileSystem{FS: fs},
	}

	v1APIGenerator.WriteWellKnownFile(ctx)

	modules, err := module.ListModules()
	if err != nil {
		slog.Error("Failed to list modules", slog.Any("err", err))
		os.Exit(1)
	}

	for _, m := range modules {
		slog.Info("Generating", slog.String("module", m.Namespace+"/"+m.Name+"/"+m.System))
		err := v1APIGenerator.GenerateModuleResponses(ctx, m.Namespace, m.Name, m.System)
		if err != nil {
			slog.Error("Failed to generate module version listing response", slog.Any("err", err))
			os.Exit(1)
		}
	}

	slog.Info("Completed generating v1 API responses")
}
