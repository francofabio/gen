package platform

import (
	"os"
	"path/filepath"
)

// ConfigDir returns the user config directory for gen.
// Linux/macOS: ~/.gen
// Windows: %USERPROFILE%\.gen
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gen"), nil
}

// ConfigFilePath returns the path to the config file.
// It does not check if the file or directory exists.
func ConfigFilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}
