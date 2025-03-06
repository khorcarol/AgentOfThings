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

	fpath := filepath.Join("images", friend.User.UserID.String())

	out, err2 := os.Create(filepath.Join(sdir, fpath))
	if err2 != nil {
		log.Fatalf("Could not create image file: %s", err2)
	}
	png.Encode(out, friend.Photo)

	out.Close()


	return FriendJson{User: friend.User, Photo: fpath, Name:friend.Name}
}

func friendJsonToFriend(fj FriendJson) api.Friend {
	sdir, err := GetStorageDir()
	if err != nil {
		log.Fatalf("Could not find storage directory: %s", err)
	}
	
	img, err2 := openImage(filepath.Join(sdir, fj.Photo))
	if err2 != nil {
		log.Fatalf("Could not open image file: %s", err2)
	}

	return api.Friend{User: fj.User, Photo: img, Name: fj.Name}
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
