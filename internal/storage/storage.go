package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
)

const (
	appDirName      = "agentofthings"
	friendsFileName = "friends.json"
)

// dirProvider interface for getting system directories
type dirProvider interface {
	GetHomeDir() (string, error)
	GetConfigDir() (string, error)
}

// defaultDirProvider implements dirProvider
type defaultDirProvider struct{}

func (p defaultDirProvider) GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (p defaultDirProvider) GetConfigDir() (string, error) {
	return os.UserConfigDir()
}

// concrete implementation of dirProvider
var provider dirProvider = defaultDirProvider{}

// returns directory where data should be stored
func GetStorageDir() (string, error) {
	configDir, err := provider.GetConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	storageDir := filepath.Join(configDir, appDirName)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	return storageDir, nil
}

// serializes and writes friends map to disk
func SaveFriends(friends map[api.ID]api.User) error {
	storageDir, err := GetStorageDir()
	if err != nil {
		return fmt.Errorf("failed to get storage directory: %w", err)
	}

	data, err := json.MarshalIndent(friends, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal friends data: %w", err)
	}

	filePath := filepath.Join(storageDir, friendsFileName)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write friends data to file: %w", err)
	}

	return nil
}

// reads and deserializes friends map from disk
// returns an empty map if the file does not exist yet
func LoadFriends() (map[api.ID]api.User, error) {
	storageDir, err := GetStorageDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage directory: %w", err)
	}

	filePath := filepath.Join(storageDir, friendsFileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[api.ID]api.User), nil
		}
		return nil, fmt.Errorf("failed to read friends data from file: %w", err)
	}

	var friends map[api.ID]api.User
	if err := json.Unmarshal(data, &friends); err != nil {
		return nil, fmt.Errorf("failed to unmarshal friends data: %w", err)
	}

	return friends, nil
}
