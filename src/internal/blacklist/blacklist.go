package blacklist

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
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

var (
	instance *Blacklist
	once     sync.Once
	mu       sync.RWMutex
)

// Load reads the blacklist from the versions_blacklist.json file
func Load() (*Blacklist, error) {
	var loadErr error
	once.Do(func() {
		// Try to find the blacklist file by up a couple directories
		// we're not always sure what the current working directory is
		possiblePaths := []string{
			"versions_blacklist.json",
			"../versions_blacklist.json",
			"../../versions_blacklist.json",
		}

		var data []byte
		var err error
		for _, path := range possiblePaths {
			data, err = os.ReadFile(path)
			if err == nil {
				break
			}
		}

		if err != nil {
			// If file doesn't exist, return empty blacklist
			if os.IsNotExist(err) {
				instance = &Blacklist{
					Providers: []ProviderEntry{},
					Modules:   []ModuleEntry{},
				}
				return
			}
			loadErr = fmt.Errorf("failed to read blacklist file: %w", err)
			return
		}

		instance = &Blacklist{}
		if err := json.Unmarshal(data, instance); err != nil {
			loadErr = fmt.Errorf("failed to parse blacklist file: %w", err)
			return
		}
	})

	return instance, loadErr
}

// IsProviderVersionBlacklisted checks if a specific provider version is blacklisted
func IsProviderVersionBlacklisted(namespace, name, version string) (bool, string) {
	mu.RLock()
	defer mu.RUnlock()

	blacklist, err := Load()
	if err != nil {
		// If there's an error loading blacklist, allow all versions
		return false, ""
	}

	for _, entry := range blacklist.Providers {
		if entry.Namespace == namespace && entry.Name == name && entry.Version == version {
			return true, entry.Reason
		}
	}
	return false, ""
}

// IsModuleVersionBlacklisted checks if a specific module version is blacklisted
func IsModuleVersionBlacklisted(namespace, name, targetSystem, version string) (bool, string) {
	mu.RLock()
	defer mu.RUnlock()

	blacklist, err := Load()
	if err != nil {
		// If there's an error loading blacklist, allow all versions
		return false, ""
	}

	for _, entry := range blacklist.Modules {
		// For modules, we need to match namespace/name/target_system/version
		if entry.Namespace == namespace && entry.Name == name && entry.TargetSystem == targetSystem && entry.Version == version {
			return true, entry.Reason
		}
	}
	return false, ""
}
