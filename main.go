package main

import (
	"image"
	"io"
	"log"
	"os"

	"github.com/khorcarol/AgentOfThings/frontend"
	hub_connection "github.com/khorcarol/AgentOfThings/hub"
	"github.com/khorcarol/AgentOfThings/internal/api"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/middle"
	"github.com/khorcarol/AgentOfThings/internal/personal"
	"github.com/khorcarol/AgentOfThings/internal/storage"
)

func main() {
	var isUser bool
	var hubName string
	isUser = true

	for i, arg := range os.Args {
		if arg == "--profile" {
			if len(os.Args) > i+1 {
				storage.SetProfileSubdirectory(os.Args[i+1])
			} else {
				log.Fatalf("%s must have a value", arg)
			}
		}
		if arg == "--hub" {
			isUser = false
			if len(os.Args) > i+1 {
				hubName = os.Args[i+1]
				storage.SetProfileSubdirectory(os.Args[i+1])
			} else {
				log.Fatalf("%s must have a value", arg)
			}
		}
	}

	if isUser {
		personal.Init()
		frontend.Init()

		if personal.IsNewUser() {
			frontend.InitLoginForm(func(name, interest, contact string, profileImageReader io.ReadCloser) {
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
	} else {
		uuid, _ := personal.GetUUID()
		hub := api.Hub{
			HubID:   api.ID{Address: uuid},
			HubName: hubName,
		}
		hub_connection.InitConnectionManager(hub)
		for {
		}
	}
}
