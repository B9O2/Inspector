package Inspect

import (
	"os"
	"path/filepath"
)

func ExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(ex)
}
