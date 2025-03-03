package personal

import (
	"image"
	"os"
	"path/filepath"
)

func GetPicture() image.Image{
	cachePath, err := os.UserCacheDir()
	path := filepath.Join(cachePath, "AgentOfThings", "profile", "profilePicture.png")

	reader, err := os.Open(path)
	if err != nil{
		// If no profile picture can be found, use the blank one
		path = "assets/blank-profile.png"
	}

	img, _, err := image.Decode(reader)

	return img 
}

// TODO: Load name
func GetName() string {
	return ""
}
