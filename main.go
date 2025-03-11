package main

import (
	"image"
	"io"
	"log"
	"os"

	"github.com/khorcarol/AgentOfThings/frontend"
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/middle"
	"github.com/khorcarol/AgentOfThings/internal/personal"
	"github.com/khorcarol/AgentOfThings/internal/storage"
)

func main() {
	for i, arg := range os.Args {
		if arg == "--profile" {
			if len(os.Args) > i+1 {
				storage.SetProfileSubdirectory(os.Args[i+1])
			} else {
				log.Fatalf("%s must have a value", arg)
			}
		}
	}

	personal.Init()
	frontend.Init()

	if personal.IsNewUser() {
		frontend.InitLoginForm(func(name, interest, contact string, profileImageReader io.ReadCloser) {
			log.Println("Ok from inside callback")
			personal.AddInterest(api.Interest{Category: 4, Description: interest})
			personal.SetPersonal(name, contact)
			connection_manager := connection.GetCMGR()
			middle.Start()
			connection_manager.StartDiscovery()

			if profileImageReader == nil {
				return
			}

			profileImage, _, err := image.Decode(profileImageReader)
			defer profileImageReader.Close()

			if err != nil {
				log.Printf("Failed to read profile image: %v", err)
			} else {
				personal.SetPicture(profileImage)
			}
		})
	} else {
		connection_manager := connection.GetCMGR()
		middle.Start()
		connection_manager.StartDiscovery()
	}

	frontend.Run()
}
