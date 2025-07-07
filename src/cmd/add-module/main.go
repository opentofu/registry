package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/opentofu/registry-stable/internal/blacklist"
	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/github"
	"github.com/opentofu/registry-stable/internal/module"

	regaddr "github.com/opentofu/registry-address"
)

type Output struct {
	File       string `json:"file"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
	Target     string `json:"target"`
	Validation string `json:"validation"`
	Exists     bool   `json:"exists"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	repository := flag.String("repository", "", "The module repository to add")
	outputFile := flag.String("output", "", "Path to write JSON result to")
	moduleDataDir := flag.String("module-data", "../modules", "Directory containing the module data")

	flag.Parse()

	ctx := context.Background()

	bl, err := blacklist.Load()
	if err != nil {
		logger.Error("Failed to load blacklist", slog.Any("err", err))
		os.Exit(1)
	} else {
		logger.Info("Loaded blacklist successfully")
	}

	token, err := github.EnvAuthToken()
	if err != nil {
		logger.Error("Initialization Error", slog.Any("err", err))
		os.Exit(1)
	}
	ghClient := github.NewClient(ctx, logger, token)

	output := Output{Exists: false}

	err = func() error {
		// Lower case input
		re := regexp.MustCompile("(?P<Namespace>[a-zA-Z0-9-]+)/terraform-(?P<Target>[a-zA-Z0-9]*)-(?P<Name>[a-zA-Z0-9-]*)")
		match := re.FindStringSubmatch(*repository)
		if match == nil {
			return fmt.Errorf("Invalid repository name: %s", *repository)
		}

		submitted := module.Module{
			Namespace:    match[re.SubexpIndex("Namespace")],
			Name:         match[re.SubexpIndex("Name")],
			TargetSystem: match[re.SubexpIndex("Target")],
			Directory:    *moduleDataDir,
			Logger:       logger,
			Github:       ghClient,
			Blacklist:    bl,
		}

		_, err = regaddr.ParseModuleSource(fmt.Sprintf("%s/%s/%s", submitted.Namespace, submitted.Name, submitted.TargetSystem))
		if err != nil {
			return err
		}

		modules, err := module.ListModules(*moduleDataDir, "", logger, ghClient, bl)
		if err != nil {
			return err
		}
		for _, p := range modules {
			if strings.ToLower(p.RepositoryURL()) == strings.ToLower(submitted.RepositoryURL()) {
				output.Exists = true
				return fmt.Errorf("Repository already exists in the registry, %s", p.RepositoryURL())
			}
		}

		err = submitted.WriteMetadata(module.Metadata{})
		if err != nil {
			return fmt.Errorf("An unexpected error occured: %w", err)
		}

		err = submitted.UpdateMetadataFile()
		if err != nil {
			return fmt.Errorf("An unexpected error occured: %w", err)
		}

		meta, err := submitted.ReadMetadata()
		if err != nil {
			return fmt.Errorf("An unexpected error occured: %w", err)
		}
		if len(meta.Versions) == 0 {
			return fmt.Errorf("No versions detected for repository %s", submitted.RepositoryURL())
		}

		output.Namespace = submitted.Namespace
		output.Name = submitted.Name
		output.Target = submitted.TargetSystem
		output.File = submitted.MetadataPath()
		return nil
	}()

	if err != nil {
		logger.Error("Unable to add module", slog.Any("err", err))
		output.Validation = err.Error()
		// Don't exit yet, still need to write the json.
	}

	jsonErr := files.SafeWriteObjectToJSONFile(*outputFile, output)
	if jsonErr != nil {
		// This really should not happen
		panic(jsonErr)
	}

	if err != nil {
		os.Exit(1)
	}
}
