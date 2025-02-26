package sources

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
)

func getSourceCacheFileName(sourceName string) (string, error) {
	cachePath, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(cachePath, "AgentOfThings", "sources", sourceName+".json")

	return path, nil
}

func cacheSourceInterests(interests []api.Interest, sourceName string) error {
	path, err := getSourceCacheFileName(sourceName)
	if err != nil {
		return err
	}

	data, err := json.Marshal(interests)
	if err != nil {
		return err
	}

	_, err = os.Stat(filepath.Dir(path))
	if errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getCachedSourceInterests(sourceName string) *[]api.Interest {
	path, err := getSourceCacheFileName(sourceName)
	if err != nil {
		return nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var interests []api.Interest

	err = json.Unmarshal(bytes, &interests)
	if err != nil {
		return nil
	}

	return &interests
}
