package main

import (
	"os"
	"path/filepath"
	"strings"
)

func isValidDir(path string) bool {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
