package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/opentofu/registry-stable/internal/module"
	"github.com/opentofu/registry-stable/internal/provider"
	"github.com/opentofu/registry-stable/internal/validate"
)

func main() {
	var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	const (
		helpStr = `Run cmd/validate <CMD> PATH/TO/definition.json

CMD:
- module: validate module's registry JSON file.
- provider: validate provider's registry JSON file.
`
		cliError = "command-line arguments don't match CLI signature"
	)

	args := os.Args[1:]

	if len(args) < 1 {
		logger.Error(cliError, slog.String("type", errorCLI))
		os.Exit(1)
	}

	if args[0] == "help" {
		fmt.Print(helpStr)
		os.Exit(0)
	}

	if len(args) != 2 {
		fmt.Print(helpStr)
		logger.Error(cliError, slog.String("type", errorCLI))
		os.Exit(1)
	}

	path := args[1]

	var (
		err     error
		errType string = errorJSONParsing
	)
	switch cmd := args[0]; cmd {
	case "module":
		var v module.Metadata
		err = readJSONFile(path, &v)
		if err == nil {
			errType = errorValidation
			err = module.Validate(v)
		}

	case "provider":
		var v provider.Metadata
		err = readJSONFile(path, &v)
		if err == nil {
			errType = errorValidation
			err = provider.Validate(v)
		}

	default:
		fmt.Print(helpStr)
		logger.Error(fmt.Sprintf("%s command is not supported", cmd), slog.String("type", errorCLI))
		os.Exit(1)
	}

	if err != nil {
		args := []any{
			slog.String("type", errType),
			slog.String("path", path),
		}
		switch err.(type) {
		case validate.Errors:
			for _, e := range err.(validate.Errors) {
				logger.Error(e.Error(), args...)
			}
		default:
			logger.Error(err.Error(), args...)
		}

		os.Exit(1)
	}
}

func readJSONFile(p string, v any) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(v)
}

const (
	errorValidation  = "validation"
	errorJSONParsing = "parsing"
	errorCLI         = "CLI"
)
