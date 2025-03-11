package personal

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/sources"
	"github.com/khorcarol/AgentOfThings/internal/storage"
)

var self api.Friend

const profilePictureFileName = "profilePicture.png"

func Init() {
	uuid, _ := GetUUID()
	id := api.ID{Address: uuid}
	interests := sources.GetInterests()
	us := api.User{UserID: id, Interests: interests, Seen: false}

	self = api.Friend{User: us, Photo: api.ImageData{Img: GetPicture()}, Name: getName()}
}

// Returns self, the Friend struct containing our personal data
func GetSelf() api.Friend {
	return self
}

func getProfilePicturePath() (string, error) {
	storagePath, err := storage.GetStorageDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(storagePath, profilePictureFileName), nil
}

func getCandidatePaths() []string {
	var paths []string
	if profilePicturePath, err := getProfilePicturePath(); err == nil {
		paths = append(paths, profilePicturePath)
	}
	paths = append(paths, filepath.Join("assets", "blank-profile.png"))
	return paths
}

func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// Iterate through possible paths to find a profile picture.
// Paths come from [getCandidatePaths], just to separate error handling neatly.
func GetPicture() image.Image {
	candidates := getCandidatePaths()
	var lastErr error
	for _, path := range candidates {
		if img, err := openImage(path); err == nil {
			return img
		} else {
			lastErr = err
		}
	}
	// No images could be loaded at all, even the basic.
	log.Print(errors.New("failed to load a profile picture: " + lastErr.Error()))
	return nil
}

func SetPicture(picture image.Image) error {
	self.Photo.Img = picture

	filePath, err := getProfilePicturePath()
	if err != nil {
		return fmt.Errorf("failed to get profile picture path: %w", err)
	}

	out, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("failed to create user image file: %w", err)
	}

	if err := png.Encode(out, picture); err != nil {
		return fmt.Errorf("failed to encode user image: %w", err)
	}

	return nil
}

func getName() string {
	name, err := storage.LoadUserName()

	if err != nil {
		name = "John Doe"
	}

	return name
}

func AddInterest(interest api.Interest) {
	if err := sources.AddManualInterest(interest); err != nil {
		log.Printf("Failed to add interest: %v\n", err)
		return
	}
	self.User.Interests = append(self.User.Interests, interest)
}

func SetName(name string) {
	storage.SaveUserName(name)
	self.Name = name
}

func SetContact(contact string) {
	self.Contact = contact
}
