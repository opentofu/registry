package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/gpg"
	"github.com/opentofu/registry-stable/internal/provider"
)

const templateString = `
terraform {
	required_providers {
		%s = {
			source = "%s/%s"
			version = "%s"
		}
	}
}
`

// cat ../log.txt  | jq -r 'select(.level = "error") | select(.type = "provider" | .out) | .out' | grep "Error while installing" | sort
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Starting terrible provider key validation")

	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")
	providerNamespace := flag.String("provider-namespace", "", "Which provider namespace to limit the command to")
	keyDataDir := flag.String("key-data", "../keys", "Directory containing the gpg keys")

	flag.Parse()

	ctx := context.Background()
	ghClient := github.NewClient(ctx, logger, "")

	providers, err := provider.ListProviders(*providerDataDir, *providerNamespace, logger, ghClient)
	if err != nil {
		logger.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}
	err = providers.Parallel(10, func(p provider.Provider) error {
		if p.Namespace == "opentofu" || p.Namespace == "hashicorp" {
			// Skip!
			return nil
		}

		keyCollection := gpg.KeyCollection{
			Namespace: p.EffectiveNamespace(),
			Directory: *keyDataDir,
		}

		keys, err := keyCollection.ListKeys()
		if err != nil {
			return err
		}

		if len(keys) == 0 {
			return nil
		}

		metadata, err := p.ReadMetadata()
		if err != nil {
			return err
		}

		pke := ProviderKeyCheckError{p: p, errs: make(map[string]error)}

		for _, ver := range metadata.Versions {

			dir, err := os.MkdirTemp("", "provider-keys")
			if err != nil {
				return err
			}
			defer os.RemoveAll(dir)

			contents := fmt.Sprintf(templateString, p.ProviderName, p.Namespace, p.ProviderName, ver.Version)
			err = os.WriteFile(dir+"/main.tf", []byte(contents), 0644)
			if err != nil {
				return err
			}

			p.Logger.Info(fmt.Sprintf("Checking version %s", ver.Version))
			cmd := exec.Command("/home/cmesh/go/bin/tofu", "init", "-no-color")
			cmd.Dir = dir
			out := new(strings.Builder)
			cmd.Stdout = out
			cmd.Stderr = out
			if err := cmd.Run(); err != nil {
				p.Logger.Info(fmt.Sprintf("Version %s failed!", ver.Version), "out", out.String())
				pke.errs[ver.Version] = fmt.Errorf("%w: %s", err, out.String())
			}
		}

		if len(pke.errs) != 0 {
			return pke
		}

		return nil
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Completed ")
}

type ProviderKeyCheckError struct {
	p    provider.Provider
	errs map[string]error
}

func (err ProviderKeyCheckError) Error() string {
	msg := fmt.Sprintf("%s/%s encountered %d versions with errors:", err.p.Namespace, err.p.ProviderName, len(err.errs))
	for ver, verr := range err.errs {
		msg += fmt.Sprintf("\n\t%s: %s", ver, verr)
	}
	return msg
}
