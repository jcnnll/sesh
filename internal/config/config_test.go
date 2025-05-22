package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jcnnll/sesh/internal/config"
)

func withTempHome(t *testing.T, testFunc func(tmpHome string)) {
	t.Helper()

	tmpDir := t.TempDir()
	config.SetUserHomeDir(func() (string, error) {
		return tmpDir, nil
	})

	testFunc(tmpDir)
}

func TestRemoveAddAndGetPaths(t *testing.T) {
	withTempHome(t, func(tmpHome string) {
		testPath := filepath.Join(tmpHome, "projects")
		os.MkdirAll(testPath, 0755)

		// add path
		err := config.AddPath("~/projects")
		if err != nil {
			t.Fatalf("AddPath failed: %v", err)
		}

		// remove path
		err = config.RemovePath("~/")
		if err != nil {
			t.Fatalf("RemovePath failed: %v", err)
		}

		paths, err := config.GetPaths()
		if err != nil {
			t.Fatalf("GetPaths failed: %v", err)
		}
		if len(paths) != 1 || paths[0] != testPath {
			t.Errorf("expected path %q, got %v", testPath, paths)
		}

		// check default editor
		ed, err := config.GetEditor()
		if err != nil {
			t.Fatalf("GetEditor failed: %v", err)
		}
		if ed != "nvim" {
			t.Errorf("expected editor 'nvim', got %v", ed)
		}

		// change editor to vim
		config.SetEditor("vim")

		ed, err = config.GetEditor()
		if err != nil {
			t.Fatalf("GetEditor failed: %v", err)
		}
		if ed != "vim" {
			t.Errorf("expected editor 'vim', got %v", ed)
		}

		// Check file contents
		cfgFile := filepath.Join(tmpHome, ".config", "sesh", "config.json")
		data, err := os.ReadFile(cfgFile)
		if err != nil {
			t.Fatalf("failed to read config file: %v", err)
		}
		if !strings.Contains(string(data), "projects") {
			t.Errorf("config file does not contain expected path: %s", data)
		}
		if !strings.Contains(string(data), "vim") {
			t.Errorf("config file does not contain expected editor: %s", data)
		}
	})
}
