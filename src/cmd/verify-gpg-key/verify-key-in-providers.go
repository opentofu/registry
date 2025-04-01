package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/opentofu/libregistry/metadata"
	"github.com/opentofu/libregistry/provider_key_verifier"
	"github.com/opentofu/libregistry/types/provider"
)

func buildKeyVerifier(keyData []byte, dataAPI metadata.API) (provider_key_verifier.ProviderKeyVerifier, error) {
	keyVerification, err := provider_key_verifier.New(keyData, dataAPI)
	if err != nil {
		return nil, err
	}
	return keyVerification, nil
}

func getProvider(ctx context.Context, dataAPI metadata.API, namespace string, providerName string) (provider.Addr, error) {
	providerAddr := provider.Addr{
		Namespace: namespace,
		Name:      providerName,
	}
	canonicalAddr, err := dataAPI.GetProviderCanonicalAddr(ctx, providerAddr)
	if err != nil {
		return provider.Addr{}, err
	}

	return canonicalAddr, nil
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

func getProviders(ctx context.Context, logger slog.Logger, dataAPI metadata.API, namespace string, providerName string) ([]provider.Addr, error) {
	var err error
	var providers []provider.Addr
	var providerAddr provider.Addr

	if len(providerName) > 0 {
		logger.Debug("Using provider", slog.String("provider", providerName))
		providerAddr, err = getProvider(ctx, dataAPI, namespace, providerName)
		providers = append(providers, providerAddr)
	} else {
		logger.Debug("Using namespace", slog.String("namespace", namespace))
		providers, err = listProviders(ctx, dataAPI, namespace)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}
	return providers, nil
}
