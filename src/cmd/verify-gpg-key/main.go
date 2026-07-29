package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/mail"
	"os"
	"regexp"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/gopenpgp/v2/crypto"

	"github.com/opentofu/registry-stable/internal/blacklist"
	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/gpg"
	"github.com/opentofu/registry-stable/internal/parallel"
	"github.com/opentofu/registry-stable/internal/provider"
	"github.com/opentofu/registry-stable/pkg/verification"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	keyFile := flag.String("key-file", "", "Location of the GPG key to verify")
	username := flag.String("username", "", "Github username to verify the GPG key against")
	orgName := flag.String("org", "", "Github organization name to verify the GPG key against")
	providerName := flag.String("provider-name", "", "Key used to sign provider-scoped name")

	outputFile := flag.String("output", "", "Path to write JSON result to")
	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")
	flag.Parse()

	logger = logger.With(slog.String("github", *username), slog.String("org", *orgName))
	slog.SetDefault(logger)
	logger.Debug("Verifying GPG key from location", slog.String("location", *keyFile))

	token, err := github.EnvAuthToken()
	if err != nil {
		logger.Error("Initialization Error", slog.Any("err", err))
		os.Exit(1)
	}

	ctx := context.Background()
	ghClient := github.NewClient(ctx, logger, token)
	ctxVerifier, cancelVerifierFn := context.WithCancel(context.Background())
	ghVerifierClient := github.NewClient(ctxVerifier, logger, token)
	bl, err := blacklist.Load()
	if err != nil {
		logger.Error("Failed to load blacklist", slog.Any("err", err))
		os.Exit(1)
	} else {
		logger.Info("Loaded blacklist successfully")
	}

	providers, err := provider.ListProviders(*providerDataDir, *orgName, logger, ghVerifierClient, bl)
	if err != nil {
		logger.Error("Failed to list providers", slog.Any("err", err))
		os.Exit(1)
	}
	var filteredProviders provider.List
	if *providerName == "" {
		filteredProviders = providers
	} else {
		for _, provider := range providers {
			if strings.EqualFold(provider.ProviderName, *providerName) {
				filteredProviders = append(filteredProviders, provider)
			}
		}
	}

	result := &verification.Result{}

	s := VerifyKey(*keyFile, filteredProviders, cancelVerifierFn)
	result.Steps = append(result.Steps, s)

	s = verification.VerifyGithubUser(ghClient, *username, *orgName)
	result.Steps = append(result.Steps, s)

	fmt.Println(result.RenderMarkdown())

	if *outputFile != "" {
		jsonErr := files.SafeWriteObjectToJSONFile(*outputFile, result.RenderMarkdown())
		if jsonErr != nil {
			// This really should not happen
			panic(jsonErr)
		}
	}

	if result.DidFail() {
		os.Exit(-1)
	}
}

var gpgNameEmailRegex = regexp.MustCompile(`.*\<(.*)\>`)

func VerifyKey(location string, providers provider.List, cancelVerifierFn context.CancelFunc) *verification.Step {
	verifyStep := &verification.Step{
		Name: "Validate GPG key",
	}

	// read the key from the filesystem
	keyData, err := os.ReadFile(location)
	if err != nil {
		verifyStep.AddError(fmt.Errorf("failed to read key file: %w", err))
		verifyStep.Status = verification.StatusFailure
		return verifyStep
	}

	var key *crypto.Key
	verifyStep.RunStep("Key is a valid PGP key", func() error {
		// From internal/gpg/key.go
		k, err := crypto.NewKeyFromArmored(string(keyData))
		if err != nil {
			return fmt.Errorf("could not parse key: %w", err)
		}
		key = k
		return nil
	})

	if key == nil {
		// The previous step failed.
		return verifyStep
	}

	expiredStep := verifyStep.RunStep("Key is not expired", func() error {
		if key.IsExpired() {
			return fmt.Errorf("key has expired, whilst this is accepted by the registry, please ensure a renewed key is used for signing new versions of providers")
		}
		return nil
	})

	// Expired keys are still accepted. They remain valid for verifying
	// signatures they produced while valid. This should be a warning and not an error
	expiredStep.FailureToWarning()

	verifyStep.RunStep("Key is not revoked", func() error {
		if key.IsRevoked() {
			return fmt.Errorf("key is revoked")
		}
		return nil
	})

	verifyStep.RunStep("Key can be used for signing", func() error {
		// Use gpg.CanSign rather than key.CanVerify so that an expired key,
		// is not rejected here purely for being expired.
		if !gpg.CanSign(key) && !key.CanVerify() {
			return fmt.Errorf("key cannot be used for signing")
		}
		return nil
	})

	emailStep := verifyStep.RunStep("Key has a valid identity and email. (Email is preferable but optional)", func() error {
		if key.GetFingerprint() == "" {
			return fmt.Errorf("key has no fingerprint")
		}

		entity := key.GetEntity()
		if entity == nil {
			return fmt.Errorf("key has no entity")
		}

		identities := entity.Identities
		if len(identities) == 0 {
			return fmt.Errorf("key has no identities")
		}

		for idName, identity := range identities {
			if identity.Name == "" {
				return fmt.Errorf("key identity %s has no name", idName)
			}

			email := gpgNameEmailRegex.FindStringSubmatch(identity.Name)
			if len(email) != 2 {
				return fmt.Errorf("key identity %s has no email", idName)
			}

			_, err := mail.ParseAddress(email[1])
			if err != nil {
				return fmt.Errorf("key identity %s has an invalid email: %w", idName, err)
			}
		}

		return nil
	})

	emailStep.FailureToWarning()

	if !verifyStep.DidFail() {
		gpgStep := verifyStep.RunStep("Key is used to sign at least one provider", func() error {
			// Inspired by OpenTofu's getproviders
			if _, err := openpgp.ReadArmoredKeyRing(strings.NewReader(string(keyData))); err != nil {
				return fmt.Errorf("error decoding signing key: %w", err)
			}

			foundProviderForKey := false

			err := providers.Parallel(20, func(p provider.Provider) error {
				meta, err := p.ReadMetadata()
				if err != nil {
					return err
				}
				meta.Logger.Info("Starting key signature checks")

				var versionChecks []parallel.Action
				for _, version := range meta.Versions {
					versionChecks = append(versionChecks, func() error {
						logger := meta.Logger.With(slog.String("version", version.Version))
						logger.Info("Begin version check")

						// Inspired by OpenTofu's getproviders
						shasumResp, err := p.Github.DownloadAssetContents(version.SHASumsURL)
						if err != nil {
							return err
						}

						sigResp, err := p.Github.DownloadAssetContents(version.SHASumsSignatureURL)
						if err != nil {
							return err
						}

						verified, err := gpg.VerifyDetachedSignature([]gpg.Key{{ASCIIArmor: string(keyData)}}, shasumResp, sigResp)
						if err != nil {
							return fmt.Errorf("error checking signature: %w", err)
						}
						if !verified {
							return nil
						}

						// Key might be expired, but that's allowed
						logger.Info("Key is valid for provider version")
						foundProviderForKey = true
						// Key was verified successfully, we can cancel all the parallelized requests
						cancelVerifierFn()
						return nil
					})
				}
				err = errors.Join(parallel.ForEach(versionChecks, 10)...)

				// TODO: Remove this check, I just used to test if we were skipping the errors correctly
				if errors.Is(err, context.Canceled) {
					meta.Logger.Info("Context cancel")
				}

				if err != nil && !errors.Is(err, context.Canceled) {
					meta.Logger.Error(err.Error())
				}
				return nil
			})
			if err != nil {
				return err
			}
			if !foundProviderForKey {
				return fmt.Errorf("key is not used to sign any known provider")
			}
			return nil
		})

		gpgStep.FailureToWarning()
	}

	return verifyStep
}
