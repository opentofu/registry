package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/gpg"
	"github.com/opentofu/registry-stable/internal/parallel"
	"github.com/opentofu/registry-stable/internal/provider"
	"github.com/opentofu/registry-stable/pkg/verification"
)

func VerifyKeyInProviders(logger *slog.Logger, ghClient github.Client, location string, orgName string) *verification.Step {

	verifyStep := &verification.Step{
		Name: "Verify GPG key in providers",
	}

	providerDataDir := "/Users/diogenesaherminio/workspace/opentofu/opentofu-registry/providers"
	providerNamespace := "wombelix"

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

	logger.Info(key.GetArmoredPublicKey())

	providers, err := provider.ListProviders(providerDataDir, providerNamespace, logger, ghClient)
	verifyStep.RunStep("Provider list is valid", func() error {
		if err != nil {
			return fmt.Errorf("could not read providers: %w", err)
		}
		return nil
	})

	signingKeyRing, err := gpg.BuildSigningKeyRing(key)
	verifyStep.RunStep("Can build a valid keyring", func() error {
		if err != nil {
			return fmt.Errorf("could not read build a keyring: %w", err)
		}
		return nil
	})

	logger.Info("Keyring", slog.String("keyring", signingKeyRing.FirstKeyID))
	logger.Info("Providers", slog.String("providers", providers[0].ProviderName))
	logger.Info("Providers", slog.String("providers", providers[0].Namespace))

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
			actions[i] = func() error {
				dir, err := os.MkdirTemp("", "provider-keys")
				if err != nil {
					return err
				}
				defer os.RemoveAll(dir)

				// Check ShaSum signature
				resp, err := http.Get(ver.SHASumsURL)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				respSig, err := http.Get(ver.SHASumsSignatureURL)
				if err != nil {
					return err
				}
				defer respSig.Body.Close()

				fileBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				providerData := crypto.NewPlainMessage(fileBytes)

				fileSigBytes, err := io.ReadAll(respSig.Body)
				if err != nil {
					return err
				}

				pgpSignature := crypto.NewPGPSignature(fileSigBytes)

				err = signingKeyRing.VerifyDetached(providerData, pgpSignature, crypto.GetUnixTime())
				if err != nil {
					return err
				}

				logger.Info("Verified", slog.String("provider", ver.SHASumsURL))

				// Option 2: Download versions and check in another way - need to talk to core devs
				// for _, target := range ver.Targets {
				// 	resp, err := http.Get(target.DownloadURL)
				// 	if err != nil {
				// 		return err
				// 	}
				// 	defer resp.Body.Close()

				// 	bodyBytes, err := io.ReadAll(resp.Body)
				// 	providerData := crypto.NewPlainMessage(bodyBytes)

				// 	verifyResult, err := signingKeyRing.VerifyDetached(providerData, crypto.GetUnixTime())
				// 	if err != nil {
				// 		return err
				// 	}

				// 	logger.Info("Verified", slog.String("provider", target.DownloadURL))
				// 	return nil
				// }
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
