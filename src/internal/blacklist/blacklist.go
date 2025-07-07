package blacklist

import (
	"encoding/json"
	"fmt"
	"os"
)

type ProviderEntry struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Reason    string `json:"reason"`
}

type ModuleEntry struct {
	Namespace    string `json:"namespace"`
	Name         string `json:"name"`
	TargetSystem string `json:"target_system"`
	Version      string `json:"version"`
	Reason       string `json:"reason"`
}

type Blacklist struct {
	Providers []ProviderEntry `json:"providers"`
	Modules   []ModuleEntry   `json:"modules"`
}

// loadFromFile reads the blacklist from a specific file path
func loadFromFile(filePath string) (*Blacklist, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, return empty blacklist
			return &Blacklist{
				Providers: []ProviderEntry{},
				Modules:   []ModuleEntry{},
			}, nil
		}
		return nil, fmt.Errorf("failed to read blacklist file: %w", err)
	}

	blacklist := &Blacklist{}
	if err := json.Unmarshal(data, blacklist); err != nil {
		return nil, fmt.Errorf("failed to parse blacklist file: %w", err)
	}

	return blacklist, nil
}

// Load reads the blacklist from the default location
func Load() (*Blacklist, error) {
	// The versions_blacklist.json file is at the repository root
	// When running bump-versions from GitHub Actions, the working directory is ./src
	// When running add-provider/add-module, the working directory is the repository root

	// Try both possible locations
	possiblePaths := []string{
		"versions_blacklist.json",    // When running from repo root (add-provider/add-module)
		"../versions_blacklist.json", // When running from src directory (bump-versions)
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			blacklist, err := loadFromFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to load blacklist from %s: %w", path, err)
			}
			return blacklist, nil
		}
	}

	// If file doesn't exist in either location, return empty blacklist
	return &Blacklist{
		Providers: []ProviderEntry{},
		Modules:   []ModuleEntry{},
	}, nil
}

// IsProviderVersionBlacklisted checks if a specific provider version is blacklisted
func (b *Blacklist) IsProviderVersionBlacklisted(namespace, name, version string) (bool, string) {
	for _, entry := range b.Providers {
		if entry.Namespace == namespace && entry.Name == name && entry.Version == version {
			return true, entry.Reason
		}
	}
	return false, ""
}

// IsModuleVersionBlacklisted checks if a specific module version is blacklisted
func (b *Blacklist) IsModuleVersionBlacklisted(namespace, name, targetSystem, version string) (bool, string) {
	for _, entry := range b.Modules {
		// For modules, we need to match namespace/name/target_system/version
		if entry.Namespace == namespace && entry.Name == name && entry.TargetSystem == targetSystem && entry.Version == version {
			return true, entry.Reason
		}
	}
	return false, ""
}
