package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	regaddr "github.com/opentofu/registry-address"
	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/gpg"
)

type Output struct {
	HasKeys bool   `json:"has_keys"`
	Message string `json:"message"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	namespace := flag.String("namespace", "", "The provider namespace to check")
	name := flag.String("name", "", "The provider name to check")
	outputFile := flag.String("output", "", "Path to write JSON result to")
	gpgDataDir := flag.String("gpg-data", "", "Directory containing the GPG key data")

	flag.Parse()

	if *namespace == "" || *name == "" || *gpgDataDir == "" {
		logger.Error("--namespace, --name, and --gpg-data are required")
		os.Exit(1)
	}

	collection := gpg.KeyCollection{
		Namespace:    *namespace,
		ProviderName: *name,
		Directory:    *gpgDataDir,
	}

	keys, err := collection.ListKeys()
	if err != nil {
		logger.Error("Failed to list GPG keys", slog.Any("err", err))
		os.Exit(1)
	}

	output := Output{HasKeys: len(keys) > 0}

	if !output.HasKeys {
		provider, providerErr := regaddr.ParseProviderSource(fmt.Sprintf("%s/%s", *namespace, *name))
		if providerErr != nil {
			logger.Error("Failed to parse provider source", slog.Any("err", providerErr))
			os.Exit(1)
		}
		output.Message = fmt.Sprintf(
			"No GPG key found for the \"%s\" provider nor for the \"%s\" namespace. "+
				"You can submit one by using this [template](https://github.com/opentofu/registry/issues/new?template=provider_key.yml). "+
				"For more details on how to do it, follow the [official guide](https://search.opentofu.org/docs/providers/adding#adding-the-gpg-key).",
			provider.String(), *namespace,
		)

		if *outputFile != "" {
			if jsonErr := files.SafeWriteObjectToJSONFile(*outputFile, output); jsonErr != nil {
				logger.Error("Failed to write output file", slog.Any("err", jsonErr))
				os.Exit(1)
			}
		}

		os.Exit(1)
	}
}
