package main

import (
	"fmt"
	"fsync/helpers"
)

func customPrint(inputStr string) {
	fmt.Printf("-------%s-------\n", inputStr)
}
func main() {
	args := helpers.ArgInit()
	switch args.Action {
	case "run":
		hosts := helpers.BuildHostConfig(args)
		customPrint("Host config built")
		hosts.VerifyHosts()
		customPrint("Hosts verified")
		hosts.StartSync()
	case "config":
		fmt.Println("Selected action config")
	}
}
