package parser

import (
	"os"
	"path/filepath"
)

func resolvePath(basePath, relativePath string) (string, error) {
	joined := filepath.Join(basePath, relativePath)

	absolutePath, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}

	return absolutePath, nil
}

func basePath(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return path, nil
	} else {
		return filepath.Dir(path), nil
	}
}
