package storage

import (
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/khorcarol/AgentOfThings/internal/api"
)

type FriendJson struct {
	User  api.User
	Photo string
	Name  string
}

func friendToFriendJson(friend api.Friend) FriendJson {
	sdir, err := GetStorageDir()
	if err != nil {
		log.Fatalf("Could not find storage directory: %s", err)
	}

	fpath := ""

	if friend.Photo.Img != nil {
		log.Print("SAVE IMG: NOT NIL")
		fpath = filepath.Join("images", friend.User.UserID.String())
		fullPath := filepath.Join(sdir, fpath)

		// Check if images folder exists, if not create it
		if _, err_img := os.Stat(filepath.Join(sdir, "images")); os.IsNotExist(err_img) {
			if err := os.Mkdir(filepath.Join(sdir, "images"), os.ModePerm); err != nil {
				log.Fatalf("Cannot create images directory: %s", err)
			}
		}

		out, err := os.Create(fullPath)
		if err != nil {
			log.Fatalf("Could not create image file: %s", err)
		}
		defer out.Close()

		if err := png.Encode(out, friend.Photo.Img); err != nil {
			log.Fatalf("Could not encode image as PNG: %s", err)
		}
	}

	return FriendJson{User: friend.User, Photo: fpath, Name: friend.Name}
}

func friendJsonToFriend(fj FriendJson) api.Friend {
	sdir, err := GetStorageDir()
	if err != nil {
		log.Fatalf("Could not find storage directory: %s", err)
	}

	var img image.Image = nil
	if fj.Photo != "" {
		img, err = openImage(filepath.Join(sdir, fj.Photo))
		if err != nil {
			log.Fatalf("Could not open image file: %s", err)
		}
	}

	return api.Friend{User: fj.User, Photo: api.ImageData{Img: img}, Name: fj.Name}
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
