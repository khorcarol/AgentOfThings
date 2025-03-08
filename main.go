package main

import (
	"log"
	"os"

	"github.com/khorcarol/AgentOfThings/frontend"
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
	frontend.Main()
}
