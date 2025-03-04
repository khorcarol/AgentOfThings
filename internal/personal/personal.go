package personal

import (
	"errors"
	"image"
	"os"
	"path/filepath"
)

func getCandidatePaths() []string {
	var paths []string
	if cachePath, err := os.UserCacheDir(); err == nil {
		paths = append(paths, filepath.Join(cachePath, "AgentOfThings", "profile", "profilePicture.png"))
	}
	paths = append(paths, "assets/blank-profile.png")
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
	panic(errors.New("failed to load a profile picture: " + lastErr.Error()))
}

// TODO: Load name.
func GetName() string {
	return "John Doe"
}
