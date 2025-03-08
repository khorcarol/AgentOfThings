package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
)

const (
	appDirName       = "AgentOfThings"
	friendsFileName  = "friends.json"
	userNameFileName = "name.txt"
)

// dirProvider interface for getting system directories
type dirProvider interface {
	GetConfigDir() (string, error)
	GetCacheDir() (string, error)
}

// defaultDirProvider implements dirProvider
type defaultDirProvider struct{}

func (p defaultDirProvider) GetConfigDir() (string, error) {
	return os.UserConfigDir()
}

func (p defaultDirProvider) GetCacheDir() (string, error) {
	return os.UserConfigDir()
}

// concrete implementation of dirProvider
var provider dirProvider = defaultDirProvider{}

var profileSubdirectory *string

func SetProfileSubdirectory(subdir string) {
	profileSubdirectory = &subdir
}

// returns directory where data should be stored
func GetStorageDir() (string, error) {
	configDir, err := provider.GetConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get config directory: %w", err)
	}

	storageDir := filepath.Join(configDir, appDirName)
	if profileSubdirectory != nil {
		storageDir = filepath.Join(storageDir, *profileSubdirectory)
	}

	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	return storageDir, nil
}

func GetCacheDir() (string, error) {
	configDir, err := provider.GetCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get cache directory: %w", err)
	}

	storageDir := filepath.Join(configDir, appDirName)
	if profileSubdirectory != nil {
		storageDir = filepath.Join(storageDir, *profileSubdirectory)
	}

	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return storageDir, nil
}

// serializes and writes friends map to disk
func SaveFriends(friends map[api.ID]api.Friend) error {

	fjm := map[api.ID]FriendJson{}
	for id, f := range friends {
		fjm[id] = friendToFriendJson(f)
	}

	storageDir, err := GetStorageDir()
	if err != nil {
		return fmt.Errorf("failed to get storage directory: %w", err)
	}

	data, err := json.MarshalIndent(fjm, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal friends data: %w", err)
	}

	filePath := filepath.Join(storageDir, friendsFileName)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write friends data to file: %w", err)
	}

	log.Printf("SAVED FILE")
	return nil
}

func SaveFriend(friend api.Friend) error {
	storageDir, err := GetStorageDir()
	if err != nil {
		return fmt.Errorf("failed to get storage directory: %w", err)
	}

	filePath := filepath.Join(storageDir, friendsFileName)
	data, err := os.ReadFile(filePath)

	fjm := make(map[api.ID]FriendJson)

	if err != nil {
		return fmt.Errorf("failed to read friends data from file: %w", err)
	} else if err := json.Unmarshal(data, &fjm); err != nil {
		return fmt.Errorf("failed to unmarshal friends data: %w", err)
	}

	fjm[friend.User.UserID] = friendToFriendJson(friend)

	wdata, err := json.MarshalIndent(fjm, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal friends data: %w", err)
	}

	if err := os.WriteFile(filePath, wdata, 0644); err != nil {
		return fmt.Errorf("failed to write friends data to file: %w", err)
	}

	return nil
}

// reads and deserializes friends map from disk
// returns an empty map if the file does not exist yet
func LoadFriends() (map[api.ID]api.Friend, error) {
	storageDir, err := GetStorageDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage directory: %w", err)
	}

	filePath := filepath.Join(storageDir, friendsFileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[api.ID]api.Friend), nil
		}
		return nil, fmt.Errorf("failed to read friends data from file: %w", err)
	}

	var fjs map[api.ID]FriendJson
	if err := json.Unmarshal(data, &fjs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal friends data: %w", err)
	}

	friends := make(map[api.ID]api.Friend)
	for id, fj := range fjs {
		friends[id] = friendJsonToFriend(fj)
	}

	return friends, nil
}

func LoadUserName() (string, error) {
	storageDir, err := GetStorageDir()
	if err != nil {
		return "", fmt.Errorf("failed to get storage directory: %w", err)
	}

	filePath := filepath.Join(storageDir, userNameFileName)
	var name []byte
	if name, err = os.ReadFile(filePath); err != nil {
		return "", fmt.Errorf("failed to read user name from file: %w", err)
	}

	return (string)(name), nil
}

func SaveUserName(name string) error {
	storageDir, err := GetStorageDir()
	if err != nil {
		return fmt.Errorf("failed to get storage directory: %w", err)
	}

	filePath := filepath.Join(storageDir, userNameFileName)
	if err := os.WriteFile(filePath, ([]byte)(name), 0644); err != nil {
		return fmt.Errorf("failed to write user name to file: %w", err)
	}

	return nil
}
