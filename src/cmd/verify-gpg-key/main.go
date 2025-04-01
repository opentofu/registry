package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/mail"
	"os"
	"regexp"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage/filesystem"
	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/gpg"
	"github.com/opentofu/registry-stable/pkg/verification"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	keyFile := flag.String("key-file", "", "Location of the GPG key to verify")
	username := flag.String("username", "", "Github username to verify the GPG key against")
	orgName := flag.String("org", "", "Github organization name to verify the GPG key against")
	providerName := flag.String("provider-name", "", "Key used to sign provider-scoped name")

	outputFile := flag.String("output", "", "Path to write JSON result to")
	providerDataDir := flag.String("provider-data", "..", "Directory containing the provider data")
	flag.Parse()

	logger = logger.With(slog.String("github", *username), slog.String("org", *orgName))
	slog.SetDefault(logger)
	logger.Debug("Verifying GPG key from location", slog.String("location", *keyFile))

	token, err := github.EnvAuthToken()
	if err != nil {
		logger.Error("Initialization Error", slog.Any("err", err))
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ghClient := github.NewClient(ctx, logger, token)

	result := &verification.Result{}

	s := VerifyKey(ctx, *logger, *providerDataDir, *keyFile, *orgName, *providerName)
	result.Steps = append(result.Steps, s)

	s = VerifyGithubUser(ghClient, *username, *orgName)
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

func VerifyGithubUser(client github.Client, username string, orgName string) *verification.Step {
	verifyStep := &verification.Step{
		Name: "Validate Github user",
	}

	s := verifyStep.RunStep(fmt.Sprintf("User is a member of the organization %s", orgName), func() error {
		member, err := client.IsUserInOrganization(username, orgName)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		if member {
			return nil
		} else {
			return fmt.Errorf("user is not a member of the organization")
		}
	})
	s.Remarks = []string{"If this is incorrect, please ensure that your organization membership is public. For more information, see [Github Docs - Publicizing or hiding organization membership](https://docs.github.com/en/account-and-profile/setting-up-and-managing-your-personal-account-on-github/managing-your-membership-in-organizations/publicizing-or-hiding-organization-membership)"}

	return verifyStep
}

var gpgNameEmailRegex = regexp.MustCompile(`.*\<(.*)\>`)

func VerifyKey(ctx context.Context, logger slog.Logger, providerDataDir string, location string, orgName string, providerName string) *verification.Step {
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
		k, err := gpg.ParseKey(string(keyData))
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

	verifyStep.RunStep("Key is not expired", func() error {
		if key.IsExpired() {
			return fmt.Errorf("key is expired")
		}
		return nil
	})

	verifyStep.RunStep("Key is not revoked", func() error {
		if key.IsRevoked() {
			return fmt.Errorf("key is revoked")
		}
		return nil
	})

	verifyStep.RunStep("Key can be used for signing", func() error {
		if !key.CanVerify() {
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

	dataAPI, err := metadata.New(filesystem.New(providerDataDir))
	if err != nil {
		verifyStep.AddError(fmt.Errorf("failed to create a metadata API: %w", err))
		verifyStep.Status = verification.StatusFailure
		return verifyStep
	}

	providers, err := getProviders(ctx, logger, dataAPI, orgName, providerName)

	if err != nil {
		verifyStep.AddError(fmt.Errorf("failed to list provider %s: %w", orgName, err))
		verifyStep.Status = verification.StatusFailure
		return verifyStep
	}

	keyVerification, err := buildKeyVerifier(keyData, dataAPI)
	if err != nil {
		verifyStep.AddError(fmt.Errorf("failed to build key verifier: %w", err))
		verifyStep.Status = verification.StatusFailure
		return verifyStep
	}

	for _, provider := range providers {
		versions, err := keyVerification.VerifyProvider(ctx, provider)
		if err != nil {
			verifyStep.AddError(fmt.Errorf("failed to verify key: %w", err))
			verifyStep.Status = verification.StatusFailure
			return verifyStep
		}

		for _, version := range versions {
			subName := fmt.Sprintf("Key is used to sign the provider %s v%s", provider, version.Version)
			verifyStep.AddStep(subName, verification.StatusSuccess)
		}
	}

	emailStep.FailureToWarning()

	return verifyStep
}
