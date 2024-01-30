package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/opentofu/registry-stable/internal/module"
	"github.com/opentofu/registry-stable/internal/provider"
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
		logger.Error(errorCLI, slog.String("error", cliError))
		os.Exit(1)
	}

	if args[0] == "help" {
		fmt.Print(helpStr)
		os.Exit(0)
	}

	if len(args) != 2 {
		fmt.Print(helpStr)
		logger.Error(errorCLI, slog.String("error", cliError))
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
		logger.Error(errorCLI, slog.String("error", fmt.Sprintf("%s command is not supported", cmd)))
		os.Exit(1)
	}

	if err != nil {
		logger.Error(errType)
		fmt.Println(err)
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
	errorValidation  = "validation error"
	errorCLI         = "CLI error"
	errorJSONParsing = "JSON parsing error"
)
