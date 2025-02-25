package middle

import (
	"log"
	"os"

	"github.com/google/uuid"
)

const uuidFilePath = ".config/uuid.txt"

// GetUUID generates a new UUID if one is not found in the persistent storage [uuidFilePath] and stores it in [uuidFilePath].
// If a UUID is found in the persistent storage, it prints the UUID to the console.
func GetUUID() uuid.UUID {
	if _, err := os.Stat(uuidFilePath); err == nil {
		// File exists.
		storedUUID, err := os.ReadFile(uuidFilePath)
		if err != nil {
			log.Fatalf("failed to read UUID from file: %v", err)
		}

		// Panics if the user has corrupted their UUID file,
		// can be stuck in a loop if the UUID file is corrupted.
		// Assume for now they don't but could be a potential issue in the future.
		return uuid.MustParse(string(storedUUID))
	} else if !os.IsNotExist(err) {
		// Shouldn't happen, if the OS plays ball.
		log.Fatalf("failed to check for UUID file: %v", err)
	}

	// File does not exist, so generate a new UUID.
	newUUID := uuid.New()
	err := os.WriteFile(uuidFilePath, []byte(newUUID.String()), 0644)
	if err != nil {
		log.Fatalf("failed to write new UUID to file: %v", err)
	}

	return newUUID
}
