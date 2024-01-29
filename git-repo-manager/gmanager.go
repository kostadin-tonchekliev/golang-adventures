package main

import (
	"fmt"
	"git-repo-manager/configHelpers"
	"os"
)

func main() {
	cliArguments := os.Args
	if len(cliArguments) > 1 {
		switch cliArguments[1] {
		case "setup":
			configHelpers.SetupEnv()
		case "config":
			configObject := configHelpers.ReadConfig()
			if len(cliArguments) > 2 {
				switch cliArguments[2] {
				case "ls":
					configObject.ListConfig()
				case "add":
					configObject.AddConfig()
				default:
					fmt.Println("[Exit 1] Please select valid sub-action")
					os.Exit(1)
				}
			} else {
				fmt.Println("[Exit 1] Please select an action")
				os.Exit(1)
			}
			configObject.CloseConfig()
		case "cd":
			configObject := configHelpers.ReadConfig()
			switch len(cliArguments) {
			case 2:
				configObject.CDRepoChoice()
			case 3:
				configObject.CDRepoManual(cliArguments[2])
			default:
				fmt.Println("[Exit 1] Please select valid sub-action")
				os.Exit(1)
			}
			configObject.CloseConfig()

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
