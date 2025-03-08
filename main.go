package main

import (
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

	middle.Start()

	connection_manager := connection.GetCMGR()
	if personal.IsNewUser() {
		frontend.InitLoginForm(func(name, interest string) {
			personal.AddInterest(api.Interest{Category: 4, Description: interest})
			personal.SetName(name)
			connection_manager.StartDiscovery()
		})
	} else {
		connection_manager.StartDiscovery()
	}

	frontend.Run()
}
