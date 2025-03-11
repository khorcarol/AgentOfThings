package hub_storage

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/storage"
)

func getStorageDir() (string, error) {
	const storagePath = "hub/messages"
	cdir, err := storage.GetCacheDir()
	if err != nil {
		return "", fmt.Errorf("Error getting cache dir: %s", err)
	}
	return filepath.Join(cdir, storagePath), nil
}

func StoreMessage(message api.Message) error {

	data, data_err := json.MarshalIndent(message, "", "  ")
	if data_err != nil {
		return data_err
	}

	sdir, path_err := getStorageDir()
	if path_err != nil {
		return path_err
	}
	path := filepath.Join(sdir, message.Timestamp.String(), message.Author.String())

	write_err := os.WriteFile(path, data, 0644)
	if write_err != nil {
		return fmt.Errorf("Error writing to file while storing message: %s", write_err)
	}

	return nil

}

func ReadMessages() ([]api.Message, error) {

	ls := []api.Message{}

	sdir, err := getStorageDir()
	if err != nil {
		return ls, err
	}

	// For every file in the path, adds the message to the list
	filepath.WalkDir(sdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		data, read_err := os.ReadFile(path)
		if read_err != nil {
			return read_err
		}

		var message api.Message
		json.Unmarshal(data, &message)

		ls = append(ls, message)

		return nil
	})

	return ls, nil

}
