package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/lib/option"
)

const (
	friendsFileName  = "friends.json"
	personalFileName = "personal.json"
)

var profileSubdirectory *string

func SetProfileSubdirectory(subdir string) {
	profileSubdirectory = &subdir
}

// returns directory where data should be stored
func GetStorageDir() (string, error) {
	storageDir, err := getAppConfigDir()
	if err != nil {
		return "", err
	}

	if profileSubdirectory != nil {
		storageDir = filepath.Join(storageDir, *profileSubdirectory)
	}

	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	return storageDir, nil
}

func GetCacheDir() (string, error) {
	cacheDir, err := getAppCacheDir()
	if err != nil {
		return "", err
	}

	if profileSubdirectory != nil {
		cacheDir = filepath.Join(cacheDir, *profileSubdirectory)
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cacheDir, nil
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
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to read friends data from file: %w", err)
		}
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

func CheckFriend(id api.ID) (option.Option[api.Friend], error) {
	storageDir, err := GetStorageDir()
	if err != nil {
		return option.OptionNil[api.Friend](), fmt.Errorf("failed to get storage directory: %w", err)
	}

	filePath := filepath.Join(storageDir, friendsFileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return option.OptionNil[api.Friend](), nil
		}
		return option.OptionNil[api.Friend](), fmt.Errorf("failed to read friends data from file: %w", err)
	}

	var fjs map[api.ID]FriendJson
	if err := json.Unmarshal(data, &fjs); err != nil {
		return option.OptionNil[api.Friend](), fmt.Errorf("failed to unmarshal friends data: %w", err)
	}

	fr, lookup := fjs[id]

	if lookup {
		return option.OptionVal(friendJsonToFriend(fr)), nil
	} else {
		return option.OptionNil[api.Friend](), nil
	}
}

type PersonalJson struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

func LoadPersonal() (PersonalJson, error) {
	dir, err := GetStorageDir()
	if err != nil {
		return PersonalJson{}, err
	}
	filePath := filepath.Join(dir, personalFileName)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return PersonalJson{Name: "John Doe", Contact: ""}, nil
		}
		return PersonalJson{}, err
	}
	var pj PersonalJson
	if err = json.Unmarshal(data, &pj); err != nil {
		return PersonalJson{}, err
	}
	return pj, nil
}

// SavePersonal writes personal details to disk.
func SavePersonal(pj PersonalJson) error {
	dir, err := GetStorageDir()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(pj, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, personalFileName), data, 0644)
}
