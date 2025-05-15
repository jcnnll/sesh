package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jcnnll/sesh/config"
)

const (
	defaultEditor = "nvim"
)

func TestGetConfig(t *testing.T) {
	tmpHome := t.TempDir()
	config.SetConfigLocation(".test_config", "test_config.json")

	config.SetUserHomeDir(func() (string, error) {
		return tmpHome, nil
	})
	defer config.SetUserHomeDir(os.UserHomeDir)

	cfg, err := config.Get()
	if err != nil {
		t.Fatalf("get config failed: %v", err)
	}

	if cfg.Editor != defaultEditor {
		t.Errorf("expected editor: %q, got: %q", defaultEditor, cfg.Editor)
	}

	if len(cfg.Paths) != 1 || cfg.Paths[0] != tmpHome {
		t.Errorf("expected paths to include only %q, got: %v", tmpHome, cfg.Paths)
	}
}

func TestSaveConfig(t *testing.T) {
	tmpHome := t.TempDir()
	config.SetConfigLocation(".test_sesh", "test_sesh.json")

	config.SetUserHomeDir(func() (string, error) {
		return tmpHome, nil
	})
	defer config.SetUserHomeDir(os.UserHomeDir)

	cfg, err := config.Get()
	if err != nil {
		t.Fatalf("get config failed: %v", err)
	}

	ed := "vim"
	cfg.Editor = ed

	cfg.AddPath("~/test")
	cfg.RemovePath(cfg.Paths[0])

	fmt.Printf("Editor:\t%s\n", cfg.Editor)
	for _, p := range cfg.Paths {
		fmt.Printf("Path:\t%s\n", p)
	}
}

func TestAddPath(t *testing.T) {
	cfg, err := config.Get()
	if err != nil {
		t.Fatalf("get config failed: %v", err)
	}

	path := "~/test"
	cfg.AddPath(path)

	if len(cfg.Paths) != 2 {
		t.Errorf("expected 2, got: %d", len(cfg.Paths))
	}

	home, _ := os.UserHomeDir()
	abs := filepath.Join(home, "test")
	if abs != cfg.Paths[1] {
		t.Errorf("expected:\n\t%v\ngot:\n\t%v", abs, cfg.Paths[1])
	}
}

func TestRemovePath(t *testing.T) {
	cfg, err := config.Get()
	if err != nil {
		t.Fatalf("get config failed: %v", err)
	}

	path := "~/test"
	cfg.AddPath(path)

	home, _ := os.UserHomeDir()
	abs := filepath.Join(home, "test")

	cfg.RemovePath(abs)
	if len(cfg.Paths) != 1 {
		t.Errorf("expected 1, got: %d", len(cfg.Paths))
	}

}

func TestSetEditor(t *testing.T) {
	cfg, err := config.Get()
	if err != nil {
		t.Fatalf("get config failed: %v", err)
	}

	ed := "vim"
	cfg.SetEditor(ed)

	if cfg.Editor != ed {
		t.Errorf("expected %s, got: %s", ed, cfg.Editor)
	}

}
