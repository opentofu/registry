package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/gpg"
	"github.com/opentofu/registry-stable/internal/parallel"
	"github.com/opentofu/registry-stable/internal/provider"
	"github.com/opentofu/registry-stable/pkg/verification"
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

func VerifyKeyInProviders(logger *slog.Logger, ghClient github.Client, location string, orgName string) *verification.Step {

	verifyStep := &verification.Step{
		Name: "Verify GPG key in providers",
	}

	providerDataDir := "../providers"
	providerNamespace := ""

	// read the key from the filesystem
	data, err := os.ReadFile(location)
	if err != nil {
		verifyStep.AddError(fmt.Errorf("failed to read key file: %w", err))
		verifyStep.Status = verification.StatusFailure
		return verifyStep
	}

	var key *crypto.Key
	verifyStep.RunStep("Key is a valid PGP key", func() error {
		k, err := gpg.ParseKey(string(data))
		if err != nil {
			return fmt.Errorf("could not parse key: %w", err)
		}
		key = k
		return nil
	})

	if key == nil {
		logger.Error("Failed to parse key", slog.Any("err", err))
	}

	providers, err := provider.ListProviders(providerDataDir, providerNamespace, logger, ghClient)
	if err != nil {
		logger.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}
	err = providers.Parallel(10, func(p provider.Provider) error {
		if p.Namespace == "opentofu" || p.Namespace == "hashicorp" {
			// Skip!
			return nil
		}

		metadata, err := p.ReadMetadata()
		if err != nil {
			return err
		}

		pke := ProviderKeyCheckError{p: p, errs: make(map[string]error)}

		actions := make([]parallel.Action, len(metadata.Versions))
		for i, ver := range metadata.Versions {
			ver := ver
			actions[i] = func() error {
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
				cmd := exec.Command("/opt/homebrew/bin/tofu", "init", "-no-color")
				cmd.Dir = dir
				out := new(strings.Builder)
				cmd.Stdout = out
				cmd.Stderr = out
				if err := cmd.Run(); err != nil {
					p.Logger.Info(fmt.Sprintf("Version %s failed!", ver.Version), "out", out.String())
					pke.errs[ver.Version] = fmt.Errorf("%w: %s", err, out.String())
				}
				return nil
			}
		}

		errs := parallel.ForEach(actions, 10)
		if len(errs) != 0 {
			return fmt.Errorf("encountered %d errors processing %d provider keys", len(errs), len(metadata.Versions))
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
	return verifyStep
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
