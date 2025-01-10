package main

import (
	"context"
	"net/http"
	"time"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/opentofu/libregistry/key_verification"
	"github.com/opentofu/libregistry/metadata/storage/filesystem"
)

func verifyKeyInProviders(ctx context.Context, providerDataDir string, key *crypto.Key, providerNamespace string) error {
	storage := filesystem.New(providerDataDir)

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	keyVerification, err := key_verification.New(httpClient, storage)
	if err != nil {
		return err
	}

	if err := keyVerification.VerifyKey(ctx, key, providerNamespace); err != nil {
		return err
	}

	return nil
}
