package main

import (
	"fmt"
	"fsync/helpers"
)

func main() {
	args := helpers.ArgInit()
	switch args.Action {
	case "run":
		fmt.Println("Selected action run")
		hosts := helpers.BuildHostConfig(args)
		hosts.VerifyHosts()
	case "config":
		fmt.Println("Selected action config")
	}
}
