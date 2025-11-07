package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/Mohammad-Ali-Rauf/sentinel.git/pkg/types"
)

// LoadConfig loads and parses the TOML configuration file
func LoadConfig(filepath string) (types.Config, error) {
	var config types.Config

	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return config, fmt.Errorf("config file not found: %s", filepath)
	}

	// Parse TOML file
	if _, err := toml.DecodeFile(filepath, &config); err != nil {
		return config, fmt.Errorf("failed to parse config: %v", err)
	}

	// Apply preset auto-fill if a known mode is set
	config.ApplyPreset()

	return config, nil
}
