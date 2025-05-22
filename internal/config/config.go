package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	Paths  []string `json:"paths"`
	Editor string   `json:"editor"`
}

var (
	defaultEditor  = "nvim"
	configDir      = ".config/sesh"
	configFileName = "config.json"
	osUserHomeDir  = os.UserHomeDir
	cfg            *config
)

func GetPaths() ([]string, error) {
	if err := load(); err != nil {
		return []string{}, err
	}

	out := make([]string, len(cfg.Paths))
	copy(out, cfg.Paths)
	return out, nil
}

func GetEditor() (string, error) {
	if err := load(); err != nil {
		return "", err
	}

	return cfg.Editor, nil
}

func SetEditor(editor string) error {
	if err := load(); err != nil {
		return err
	}

	cfg.Editor = editor
	return save()
}

func AddPath(path string) error {
	if err := load(); err != nil {
		return err
	}

	absPath, err := absPath(path)
	if err != nil {
		return fmt.Errorf("failed to add path: %w", err)
	}
	cfg.Paths = append(cfg.Paths, absPath)
	return save()
}

func RemovePath(target string) error {
	if err := load(); err != nil {
		return err
	}

	absTarget, err := absPath(target)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	updated := []string{}
	found := false

	for _, path := range cfg.Paths {
		if path == absTarget {
			found = true
			continue
		}
		updated = append(updated, path)
	}

	if !found {
		return fmt.Errorf("path not found: %s", absTarget)
	}

	cfg.Paths = updated
	return save()
}

// --- internals ---

func load() error {
	if cfg != nil {
		return nil // already loaded
	}

	homePath, err := osUserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine user home directory: %w", err)
	}

	configPath := filepath.Join(homePath, configDir, configFileName)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg = &config{
			Paths:  []string{homePath},
			Editor: defaultEditor,
		}
		return save()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	cfg = &config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	if cfg.Paths == nil {
		cfg.Paths = []string{homePath}
	}
	if cfg.Editor == "" {
		cfg.Editor = defaultEditor
	}
	return nil
}

func save() error {
	homePath, err := osUserHomeDir()
	if err != nil {
		return fmt.Errorf("could not stat user home path %w", err)
	}

	configPath := filepath.Join(homePath, configDir, configFileName)

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config file path %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config %w", err)
	}
	return os.WriteFile(configPath, data, 0644)
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
	return filepath.Abs(path)
}
