package main

import (
	"fmt"
	"os"

	"github.com/efejjota/concourse-autopilot/cmd"
)

func main() {
	autopilotCommand := os.Args[1]
	switch autopilotCommand {
	case "init":
		cmd.Init()
	case "gen":
		cmd.Gen()
	case "fly":
		cmd.Fly()
	default:
		fmt.Println("Wrong command")
	}
}
