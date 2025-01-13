package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/provider_verifier"
	"github.com/opentofu/libregistry/types/provider"
)

func buildKeyVerifier(dataAPI metadata.API) (provider_verifier.KeyVerification, error) {
	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}

	keyVerification, err := provider_verifier.New(httpClient, dataAPI)
	if err != nil {
		return nil, err
	}
	return keyVerification, nil
}

func listProviders(ctx context.Context, dataAPI metadata.API, namespace string) ([]provider.Addr, error) {
	providers, err := dataAPI.ListProvidersByNamespace(ctx, namespace, false)
	if err != nil {
		return nil, err
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("there are no providers in namespace %s; please submit at least one provider before submitting a GPG key", namespace)
	}

	return providers, nil
}
