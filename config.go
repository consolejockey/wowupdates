package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const CONFIGFILE = "config.yaml"

// Represents the underlying configuration
type Config struct {
	AddonPath string          `yaml:"addonPath"`
	Addons    map[string]bool `yaml:"addons"`
	ModIDs    map[string]int  `yaml:"modIDs"`
	API       string          `yaml:"api"`
}

// Constructor for the Config struct
func NewConfig() (Config, error) {
	file, err := os.Open(CONFIGFILE)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open YAML file: %v", err)
	}
	defer file.Close()

	yamlContent, err := io.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read YAML content: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(yamlContent, &config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	if _, err := os.Stat(config.AddonPath); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("addonPath does not exist: %v\nPlease edit the %s file", CONFIGFILE, err)
	}
	return config, nil
}

// Returns a map of addonName: addonID of enabled addons
func getEnabledAddonModIDs(config Config) map[string]int {
	enabledAddonModIDs := make(map[string]int)

	for addon, enabled := range config.Addons {
		if enabled {
			modID, exists := config.ModIDs[addon]
			if exists {
				enabledAddonModIDs[addon] = modID
			} else {
				log.Println("Oh no! ModID not found for addon:", addon)
			}
		}
	}
	return enabledAddonModIDs
}
