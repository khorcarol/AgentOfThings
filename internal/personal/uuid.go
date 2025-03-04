package personal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

const (
	appCacheDirName = "AgentOfThings"
	uuidFileName    = "uuid.json"
)

type uuidConfig struct {
	UUID string `json:"uuid"`
}

var (
	cachedUUID uuid.UUID
	initOnce   sync.Once
	initErr    error
)

// GetUUID returns the application UUID by reading it from the OS-specific
// configuration storage or generating it if it does not exist.
// The function caches the UUID so that repeated calls avoid unnecessary file I/O.
// On error, a non-nil error is returned.
func GetUUID() (uuid.UUID, error) {
	initOnce.Do(func() {
		cachedUUID, initErr = getUUIDInternal()
	})
	return cachedUUID, initErr
}

// getUUIDInternal reads the UUID from the appdata cache file, or creates a cache file if not found.
func getUUIDInternal() (uuid.UUID, error) {
	// Obtain the platform-specific configuration directory.
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user cache dir: %w", err)
	}

	appConfigPath := filepath.Join(cacheDir, appCacheDirName)
	uuidFilePath := filepath.Join(appConfigPath, uuidFileName)

	if info, err := os.Stat(uuidFilePath); err == nil && !info.IsDir() {
		// File exists; read and unmarshal the stored UUID.
		stored, err := os.ReadFile(uuidFilePath)
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to read UUID from file: %w", err)
		}

		var config uuidConfig
		if err := json.Unmarshal(stored, &config); err != nil {
			return uuid.Nil, fmt.Errorf("failed to parse UUID JSON: %w", err)
		}

		parsed, err := uuid.Parse(config.UUID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid UUID in file: %w", err)
		}

		log.Printf("UUID found in config: %s", parsed)
		return parsed, nil
	} else if err != nil && !os.IsNotExist(err) {
		// Shouldn't happen if the OS plays ball.
		return uuid.Nil, fmt.Errorf("failed to check for UUID file: %w", err)
	}

	if err := os.MkdirAll(appConfigPath, 0755); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	newUUID := uuid.New()
	config := uuidConfig{UUID: newUUID.String()}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to marshal UUID JSON: %w", err)
	}

	if err := os.WriteFile(uuidFilePath, data, 0644); err != nil {
		return uuid.Nil, fmt.Errorf("failed to write new UUID to file: %w", err)
	}

	log.Printf("Generated new UUID and saved to %s: %s", uuidFilePath, newUUID)
	return newUUID, nil
}
