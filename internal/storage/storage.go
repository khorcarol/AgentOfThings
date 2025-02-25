package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
)

const friendsFileName = "friends.json"

// returns directory where data should be stored
func GetStorageDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	storageDir := filepath.Join(homeDir, ".agentofthings")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", err
	}
	return storageDir, nil
}

// writes friends map to disk
func SaveFriends(friends map[api.ID]api.User) error {
	storageDir, err := GetStorageDir()
	if err != nil {
		return err
	}

	data, err := json.Marshal(friends)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(storageDir, friendsFileName), data, 0644)
}

// reads friends map from disk
func LoadFriends() (map[api.ID]api.User, error) {
	storageDir, err := GetStorageDir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(storageDir, friendsFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[api.ID]api.User), nil
		}
		return nil, err
	}

	var friends map[api.ID]api.User
	if err := json.Unmarshal(data, &friends); err != nil {
		return nil, err
	}

	return friends, nil
}
