package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/opentofu/registry-stable/internal/module"
	"github.com/opentofu/registry-stable/internal/provider"
	"golang.org/x/mod/semver"
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
		logger.Error(cliError)
		os.Exit(1)
	}

	if args[0] == "help" {
		fmt.Print(helpStr)
		os.Exit(0)
	}

	if len(args) != 2 {
		fmt.Print(helpStr)
		logger.Error(cliError)
		os.Exit(1)
	}

	path := args[1]

	var err error
	switch cmd := args[0]; cmd {
	case "module":
		err = validateModuleFile(path)
	case "provider":
		err = validateProviderFile(path)

	default:
		fmt.Print(helpStr)
		logger.Error("%s command is not supported", cmd)
		os.Exit(1)
	}

	if err != nil {
		logger.Error("validation error", slog.Any("error", err))
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

func validateModuleFile(p string) error {
	var v module.Metadata
	if err := readJSONFile(p, &v); err != nil {
		return err
	}

	if len(v.Versions) < 1 {
		return EmptyList
	}

	var errs = make([]error, 0, len(v.Versions))
	for _, ver := range v.Versions {
		if !isValidVersion(ver.Version) {
			err := fmt.Errorf("found semver-incompatible version: %s", ver.Version)
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func isValidVersion(ver string) bool {
	if ver[0] != 'v' {
		ver = "v" + ver
	}
	return semver.IsValid(ver)
}

func validateProviderFile(p string) error {
	var v provider.Metadata
	if err := readJSONFile(p, &v); err != nil {
		return err
	}

	if len(v.Versions) < 1 {
		return EmptyList
	}

	var errs = make([]error, 0, len(v.Versions))
	for _, ver := range v.Versions {
		if err := validateProviderVersion(ver); err != nil {
			errs = append(errs, fmt.Errorf("invalid ver %s: %w", ver.Version, err))
		}
	}

	return errors.Join(errs...)
}

func validateProviderVersion(ver provider.Version) error {
	var errs = make([]error, 0)

	if !isValidVersion(ver.Version) {
		errs = append(errs, fmt.Errorf("found semver-incompatible version: %s", ver.Version))
	}

	// validate provider's protocols:

	// TODO: add explicit validation of the protocol version
	if len(ver.Protocols) == 0 {
		errs = append(errs, fmt.Errorf("empty protocols list"))
	}

	// TODO: add validation of provider targets. Per target:
	//  - validate if os and arch are in the list of allowed
	//  - validate if the filename is consistent with the version, os and arch
	//  - validate if filename matches the url ending

	return errors.Join(errs...)
}

var (
	EmptyList = errors.New("found empty list of versions")
)
