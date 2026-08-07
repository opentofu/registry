package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/opentofu/registry-stable/internal/blacklist"
	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/gpg"
	"github.com/opentofu/registry-stable/internal/provider"
	"github.com/opentofu/registry-stable/pkg/verification"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	namespace := flag.String("namespace", "", "Provider namespace (also the GitHub org the requester must belong to)")
	name := flag.String("name", "", "Provider name")
	version := flag.String("version", "", "The provider version to re-index")
	username := flag.String("username", "", "GitHub username requesting the re-index")

	providerDataDir := flag.String("provider-data", "../providers", "Directory containing the provider data")
	gpgDataDir := flag.String("gpg-data", "../keys", "Directory containing the GPG key data")
	outputFile := flag.String("output", "", "Path to write the result to")
	flag.Parse()
	if *namespace == "" || *name == "" || *version == "" || *username == "" {
		logger.Error("--namespace, --name, --version, and --username are required")
		os.Exit(1)
	}

	logger = logger.With(slog.String("namespace", *namespace), slog.String("name", *name), slog.String("version", *version))
	slog.SetDefault(logger)

	token, err := github.EnvAuthToken()
	if err != nil {
		logger.Error("Initialization Error", slog.Any("err", err))
		os.Exit(1)
	}
	ghClient := github.NewClient(context.Background(), logger, token)

	bl, err := blacklist.Load()
	if err != nil {
		logger.Error("Failed to load blacklist", slog.Any("err", err))
		os.Exit(1)
	}

	p := provider.Provider{
		Namespace:    *namespace,
		ProviderName: *name,
		Directory:    *providerDataDir,
		Logger:       logger,
		Github:       ghClient,
		Blacklist:    bl,
	}
	result := &verification.Result{}

	result.Steps = append(result.Steps, verification.VerifyGithubUser(ghClient, *username, *namespace))

	keys, keyStep := loadKeys(*namespace, *name, *gpgDataDir)
	result.Steps = append(result.Steps, keyStep)

	result.Steps = append(result.Steps, VerifyVersionAgainstKeys(p, *version, keys))

	fmt.Println(result.RenderMarkdown())

	if *outputFile != "" {
		if jsonErr := files.SafeWriteObjectToJSONFile(*outputFile, result.RenderMarkdown()); jsonErr != nil {
			panic(jsonErr)
		}
	}
	if result.DidFail() {
		os.Exit(-1)
	}
}

func loadKeys(namespace, name, gpgDataDir string) ([]gpg.Key, *verification.Step) {
	step := &verification.Step{Name: "Validate GPG key on file"}

	collection := gpg.KeyCollection{
		Namespace:    namespace,
		ProviderName: name,
		Directory:    gpgDataDir,
	}

	var keys []gpg.Key
	step.RunStep("A GPG key has been submitted for this provider", func() error {
		k, err := collection.ListKeys()
		if err != nil {
			return fmt.Errorf("failed to list keys: %w", err)
		}
		if len(k) == 0 {
			return fmt.Errorf("no GPG key found for this provider or its namespace")
		}
		keys = k
		return nil
	})

	return keys, step
}

func VerifyVersionAgainstKeys(p provider.Provider, requestedVersion string, keys []gpg.Key) *verification.Step {
	step := &verification.Step{Name: "Validate the requested version against the GPG key"}

	step.RunStep(fmt.Sprintf("Version %s is currently indexed", requestedVersion), func() error {
		meta, err := p.ReadMetadata()
		if err != nil {
			return fmt.Errorf("failed to read provider metadata: %w", err)
		}

		for _, v := range meta.Versions {
			if v.Version == requestedVersion {
				return verifySignature(p, v, keys)
			}
		}

		return fmt.Errorf("version %q is not currently indexed for %s/%s", requestedVersion, p.Namespace, p.ProviderName)
	})

	return step
}

func verifySignature(p provider.Provider, v provider.Version, keys []gpg.Key) error {
	if len(keys) == 0 {
		// The "key on file" step already reported this
		return fmt.Errorf("no keys available to verify against")
	}

	shasums, err := p.Github.DownloadAssetContents(v.SHASumsURL)
	if err != nil {
		return fmt.Errorf("failed to download current SHASUMS: %w", err)
	}

	signature, err := p.Github.DownloadAssetContents(v.SHASumsSignatureURL)
	if err != nil {
		return fmt.Errorf("failed to download current SHASUMS signature: %w", err)
	}

	verified, err := gpg.VerifyDetachedSignature(keys, shasums, signature)
	if err != nil {
		return fmt.Errorf("signature check failed: %w", err)
	}
	if !verified {
		return fmt.Errorf("current release does not verify against any key on file for this provider")
	}

	return nil
}
