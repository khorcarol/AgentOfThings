package sources

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/storage"
)

var manualInterests []api.Interest

func getManualInterestsPath(configPath string) string {
	return filepath.Join(configPath, "manualInterests.json")
}

func loadManualInterests() []api.Interest {
	configPath, err := storage.GetStorageDir()
	if err != nil {
		log.Printf("failed to get storage directory: %v", err)
		return nil
	}

	interestsJson, err := os.ReadFile(getManualInterestsPath(configPath))

	if err != nil {
		log.Printf("failed to unmarshal interests data: %v", err)
		return nil
	}

	var interests []api.Interest
	json.Unmarshal(interestsJson, &interests)

	return interests
}

func AddManualInterest(interest api.Interest) error {
	configPath, err := storage.GetStorageDir()
	if err != nil {
		return err
	}

	updatedManualInterests := append(manualInterests, interest)
	interestsJson, err := json.Marshal(updatedManualInterests)

	if err != nil {
		return fmt.Errorf("failed to marshal interests data: %w", err)
	}

	if err := os.WriteFile(getManualInterestsPath(configPath), interestsJson, 0644); err != nil {
		return err
	}

	manualInterests = updatedManualInterests
	return nil
}
