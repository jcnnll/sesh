package sesh

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Session struct {
	Name    string
	Windows []window
}

type window struct {
	Title   string
	Command string
}

type Config struct {
	Paths  []string `json:"paths"`
	Editor string   `json:"editor"`
}

var (
	defaultEditor  = "nvim"
	configDir      = ".config/sesh"
	configFileName = "config.json"
	osUserHomeDir  = os.UserHomeDir
)

func SetConfigLocation(dir, filename string) {
	configDir = dir
	configFileName = filename
}

func (c *Config) AddPath(path string) error {
	absPath, err := absPath(path)
	if err != nil {
		return fmt.Errorf("failed to add path to config: %w", err)
	}

	c.Paths = append(c.Paths, absPath)
	return nil
}

func (c *Config) RemovePath(target string) error {
	absTarget, err := absPath(target)
	if err != nil {
		return fmt.Errorf("failed to resolve path '%s': %w", target, err)
	}

	updated := []string{}
	found := false

	for _, path := range c.Paths {
		if path == absTarget {
			found = true
			continue
		}
		updated = append(updated, path)
	}

	if !found {
		return fmt.Errorf("path not found in config: %s", err)
	}

	c.Paths = updated
	return nil
}

func (c *Config) SetEditor(editor string) {
	c.Editor = editor
}

func Get() (*Config, error) {
	homePath, err := osUserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not stat user home path %w", err)
	}

	configPath := filepath.Join(homePath, configDir, configFileName)

	// handle config file does not exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := Config{
			Paths: []string{
				homePath,
			},
			Editor: defaultEditor,
		}
		return &cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	homePath, err := osUserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir %w", err)
	}
	configFile := filepath.Join(homePath, configDir, configFileName)

	if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
		return fmt.Errorf("failed to create config file path %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config %w", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write defaut config %w", err)
	}
	return nil
}

func SetUserHomeDir(fn func() (string, error)) {
	osUserHomeDir = fn
}

func absPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := osUserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("could not resolve absolute path: %s", path)
	}
	return abs, nil
}
