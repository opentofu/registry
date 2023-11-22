package common

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"registry-stable/internal/github"
	"registry-stable/internal/module"
	"registry-stable/internal/provider"
)

type CLI struct {
	Logger          *slog.Logger
	Github          github.Client
	ModuleDataDir   string
	ProviderDataDir string
}

func Parse() CLI {
	moduleDataDir := flag.String("module-data", "../modules", "Directory containing the module data")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")

	flag.Parse()

	cli := CLI{
		Logger:          slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		ModuleDataDir:   *moduleDataDir,
		ProviderDataDir: *providerDataDir,
	}

	return cli
}

func (cli CLI) Modules(action func(m module.Module) error) error {
	modules, err := module.ListModules(cli.ModuleDataDir, cli.Logger, cli.Github)
	if err != nil {
		return fmt.Errorf("failed to list modules: %w", err)
	}

	actions := make([]func() error, len(modules))
	for i, m := range modules {
		actions[i] = func() error { return action(m) }
	}
	return cli.Run(actions)
}

func (cli CLI) Providers(action func(p provider.Provider) error) error {
	providers, err := provider.ListProviders(cli.ProviderDataDir, cli.Logger, cli.Github)
	if err != nil {
		return fmt.Errorf("failed to list providers: %w", err)
	}

	actions := make([]func() error, len(providers))
	for i, p := range providers {
		actions[i] = func() error { return action(p) }
	}
	return cli.Run(actions)
}

func (c CLI) Run(actions []func() error) error {
	errChan := make(chan error, len(actions))

	for _, a := range actions {
		a := a
		go func() {
			errChan <- a()
		}()
	}

	var errs []error
	for _ = range actions {
		err := <-errChan
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		for _, err := range errs {
			c.Logger.Error(err.Error())
		}
		return fmt.Errorf("Encountered %d errors", len(errs))
	}
	return nil
}
