package hub_storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/storage"
)

func getStorageDir() (string, error) {
	const storagePath = "hub/messages"
	cdir, err := storage.GetStorageDir()
	if err != nil {
		return "", fmt.Errorf("error getting cache dir: %s", err)
	}

	path := filepath.Join(cdir, storagePath)

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return path, nil
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

	filename := fmt.Sprintf("%v_%s", message.Timestamp.Unix(), message.Author)
	path := filepath.Join(sdir, filename)

	write_err := os.WriteFile(path, data, 0644)
	if write_err != nil {
		return fmt.Errorf("error writing to file while storing message: %s", write_err)
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

		if d.IsDir() {
			return nil
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
