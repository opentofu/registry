package v1api

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/gpg"
	"github.com/opentofu/registry-stable/internal/provider"

	cp "github.com/otiai10/copy"
)

// ProviderGenerator is responsible for generating the response for the provider version listing API endpoints.
type ProviderGenerator struct {
	provider.Provider
	provider.Metadata

	KeyLocation string
	Destination string
	log         *slog.Logger
}

// NewProviderGenerator creates a new ProviderGenerator which will generate the response for the provider version listing API endpoints and write it to the given destination.
func NewProviderGenerator(p provider.Provider, destination string, gpgKeyLocation string) (ProviderGenerator, error) {
	metadata, err := p.ReadMetadata()
	if err != nil {
		return ProviderGenerator{}, err
	}

	return ProviderGenerator{
		Provider: p,
		Metadata: metadata,

		KeyLocation: gpgKeyLocation,
		Destination: destination,
		log:         p.Logger,
	}, err
}

// VersionListingPath returns the path to the provider version listing file.
func (p ProviderGenerator) VersionListingPath() string {
	namespacePath := strings.ToLower(p.Provider.Namespace)
	providerNamePath := strings.ToLower(p.Provider.ProviderName)
	return filepath.Join(p.Destination, "v1", "providers", namespacePath, providerNamePath, "versions")
}

// VersionDownloadPath returns the path to the provider version download file.
func (p ProviderGenerator) VersionDownloadPath(ver provider.Version, details ProviderVersionDetails) string {
	namespacePath := strings.ToLower(p.Provider.Namespace)
	providerNamePath := strings.ToLower(p.Provider.ProviderName)
	version := strings.ToLower(ver.Version)
	return filepath.Join(p.Destination, "v1", "providers", namespacePath, providerNamePath, version, "download", details.OS, details.Arch)
}

// VersionListing will take the provider metadata and generate the responses for the provider version listing API endpoints.
func (p ProviderGenerator) VersionListing() ProviderVersionListingResponse {
	versions := make([]ProviderVersionResponseItem, len(p.Metadata.Versions))

	for versionIdx, ver := range p.Metadata.Versions {
		verResp := ProviderVersionResponseItem{
			Version:   strings.ToLower(ver.Version),
			Protocols: ver.Protocols,
			Platforms: make([]Platform, len(ver.Targets)),
		}

		for targetIdx, target := range ver.Targets {
			verResp.Platforms[targetIdx] = Platform{
				OS:   target.OS,
				Arch: target.Arch,
			}
		}
		versions[versionIdx] = verResp
	}

	return ProviderVersionListingResponse{
		versions,
		p.Metadata.Warnings,
	}
}

// VersionDetails will take the provider metadata and generate the responses for the provider version download API endpoints.
func (p ProviderGenerator) VersionDetails() (map[string]ProviderVersionDetails, error) {
	versionDetails := make(map[string]ProviderVersionDetails)

	keyCollection := gpg.KeyCollection{
		Namespace:    p.Provider.EffectiveNamespace(),
		ProviderName: p.Provider.ProviderName,
		Directory:    p.KeyLocation,
	}

	keys, err := keyCollection.ListKeys()
	if err != nil {
		p.log.Error("Failed to list keys", slog.Any("err", err))
		return nil, err
	}

	for _, ver := range p.Metadata.Versions {
		for _, target := range ver.Targets {
			details := ProviderVersionDetails{
				Protocols:           ver.Protocols,
				OS:                  target.OS,
				Arch:                target.Arch,
				Filename:            target.Filename,
				DownloadURL:         target.DownloadURL,
				SHASumsURL:          ver.SHASumsURL,
				SHASumsSignatureURL: ver.SHASumsSignatureURL,
				SHASum:              target.SHASum,
				SigningKeys: SigningKeys{
					GPGPublicKeys: keys,
				},
			}
			versionDetails[p.VersionDownloadPath(ver, details)] = details
		}
	}
	return versionDetails, nil
}

// Generate generates the responses for the provider version listing API endpoints.
func (p ProviderGenerator) Generate() error {
	p.log.Info("Generating")

	details, err := p.VersionDetails()
	if err != nil {
		return err
	}

	for location, details := range details {
		err := files.SafeWriteObjectToJSONFile(location, details)
		if err != nil {
			return fmt.Errorf("failed to write metadata version download file: %w", err)
		}
	}

	err = files.SafeWriteObjectToJSONFile(p.VersionListingPath(), p.VersionListing())
	if err != nil {
		return err
	}

	p.log.Info("Generated")

	return nil
}

func ArchivedOverrides(destDir string, log *slog.Logger) error {
	re := regexp.MustCompile("(?P<Namespace>.*)/terraform-provider-(?P<Name>.*)")
	namespaces := []string{"hashicorp", "opentofu"}

	for original, replacement := range provider.ArchivedOverrides {
		oMatch := re.FindStringSubmatch(strings.ToLower(original))
		rMatch := re.FindStringSubmatch(strings.ToLower(replacement))
		if oMatch == nil {
			return fmt.Errorf("invalid archived override: %s!", oMatch)
		}
		if rMatch == nil {
			return fmt.Errorf("invalid archived override: %s!", rMatch)
		}

		for _, namespace := range namespaces {
			oPath := filepath.Join(destDir, "v1", "providers", namespace, oMatch[re.SubexpIndex("Name")])
			rPath := filepath.Join(destDir, "v1", "providers", rMatch[re.SubexpIndex("Namespace")], rMatch[re.SubexpIndex("Name")])

			log.Info(fmt.Sprintf("Adding %s override from %s -> %s", namespace, rPath, oPath))

			if _, err := os.Stat(oPath); err == nil {
				return fmt.Errorf("invalid %s override: %s already exists", namespace, oPath)
			}
			if _, err := os.Stat(rPath); err != nil {
				log.Warn(fmt.Sprintf("Skipping %s override: %s does not exist", namespace, rPath))
				continue
			}
			err := cp.Copy(rPath, oPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
