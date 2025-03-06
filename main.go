package main

import (
	"log"

	"github.com/khorcarol/AgentOfThings/frontend"
	"github.com/khorcarol/AgentOfThings/internal/connection"
	"github.com/khorcarol/AgentOfThings/internal/middle"
)

func main() {
	connection_manager, err := connection.GetCMGR()
	if err != nil {
		log.Fatal("Failed to initialise ConnectionManager:", err)
	}

	frontend.Init()
	middle.Start()
	connection_manager.StartDiscovery()
	frontend.Main()
}
