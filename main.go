package main

import (
	"os"

	"github.com/khorcarol/AgentOfThings/frontend"
	"github.com/khorcarol/AgentOfThings/internal/personal"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		personal.UUIDFileName = args[0]
	}

	personal.Init()
	frontend.Init()
	frontend.Main()
}
