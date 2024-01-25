package main

import (
	"fmt"
	"git-repo-manager/configHelpers"
	"os"
)

func main() {
	cliArguments := os.Args
	if len(cliArguments) > 1 {
		// Read config in the beginning
		configObject := configHelpers.ReadConfig()

		switch cliArguments[1] {
		case "config":
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
		case "cd":
			switch len(cliArguments) {
			case 2:
				configObject.CDRepoChoice()
			case 3:
				configObject.CDRepoManual(cliArguments[2])
			default:
				fmt.Println("[Exit 1] Please select valid sub-action")
				os.Exit(1)
			}

		default:
			fmt.Println("[Exit 1] Please select valid action")
			os.Exit(1)
		}

		// Close config at the end
		configObject.CloseConfig()
	} else {
		// Can print help here
		fmt.Println("Please select an action")
		os.Exit(1)
	}
}
