//go:build !android
// +build !android

package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

const appDirName = "AgentOfThings"

func getAppCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get cache directory: %w", err)
	}

	return filepath.Join(cacheDir, appDirName), nil
}

func getAppConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	return filepath.Join(configDir, appDirName), nil
}
