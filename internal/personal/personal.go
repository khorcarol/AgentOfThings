package personal

import (
	"errors"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/sources"
	"github.com/khorcarol/AgentOfThings/internal/storage"
)

var self api.Friend

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

func getCandidatePaths() []string {
	var paths []string
	if cachePath, err := storage.GetCacheDir(); err == nil {
		paths = append(paths, filepath.Join(cachePath, "profile", "profilePicture.png"))
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
