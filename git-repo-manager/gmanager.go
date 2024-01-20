package main

import (
	"fmt"
	"git-repo-manager/step_config"
	"os"
)

func main() {
	cliArguments := os.Args
	if len(cliArguments) > 1 {
		switch cliArguments[1] {
		case "config":
			if len(cliArguments) > 2 {
				switch cliArguments[2] {
				case "ls":
					configObject := configHelpers.ReadConfig()
					fmt.Println(configObject.Name())
					configHelpers.CloseConfig(configObject)
				default:
					fmt.Println("[Exit 1] Please select valid subaction")
					os.Exit(1)
				}
			} else {
				fmt.Println("[Exit 1] Please select an action")
				os.Exit(1)
			}
		default:
			fmt.Println("[Exit 1] Please select valid action")
			os.Exit(1)
		}
	} else {
		// Can print help here
		fmt.Println("Please select an action")
		os.Exit(1)
	}
}
