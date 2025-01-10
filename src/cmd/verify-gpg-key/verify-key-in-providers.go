package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/metadata/storage"
	"github.com/opentofu/libregistry/provider_verifier"
	"github.com/opentofu/libregistry/types/provider"
)

func buildKeyVerifier(storageAPI storage.API) (provider_verifier.KeyVerification, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	keyVerification, err := provider_verifier.New(httpClient, storageAPI)
	if err != nil {
		return nil, err
	}
	return keyVerification, nil
}

func listProviders(ctx context.Context, storageAPI storage.API, namespace string) ([]provider.Addr, error) {
	dataAPI, err := metadata.New(storageAPI)
	if err != nil {
		return nil, err
	}

	providers, err := dataAPI.ListProvidersByNamespace(ctx, namespace, false)
	if err != nil {
		return nil, err
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no providers found for namespace %s", namespace)
	}

	return providers, nil
}
