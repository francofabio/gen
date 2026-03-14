package config

import (
	"encoding/json"
	"os"

	"github.com/francofabio/gen/internal/platform"
)

// Config holds the user configuration (e.g. card BINs).
type Config struct {
	Cards map[string][]string `json:"cards"`
}

// Load reads the config from ~/.gen/config.json.
// If the file does not exist, returns a non-nil Config with empty Cards (no error).
// If the file exists and is invalid JSON, returns an error.
func Load() (*Config, error) {
	path, err := platform.ConfigFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Cards: make(map[string][]string)}, nil
		}
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.Cards == nil {
		c.Cards = make(map[string][]string)
	}
	return &c, nil
}

// EnsureConfigDir creates the config directory if it does not exist.
// Used only when we need to write; not required for Load.
func EnsureConfigDir() (string, error) {
	dir, err := platform.ConfigDir()
	if err != nil {
		return "", err
	}
	return dir, os.MkdirAll(dir, 0755)
}
