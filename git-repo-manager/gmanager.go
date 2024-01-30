package main

import (
	"fmt"
	"git-repo-manager/configActions"
	"git-repo-manager/generalHelpers"
	"os"
)

func main() {
	cliArguments := os.Args
	if len(cliArguments) > 1 {
		switch cliArguments[1] {
		case "setup":
			configActions.SetupEnv()
		case "config":
			generalHelpers.VerifyEnv()
			configObject := configActions.ReadConfig()
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
			configObject.CloseFiles()
		case "cd":
			generalHelpers.VerifyEnv()
			configObject := configActions.ReadConfig()
			switch len(cliArguments) {
			case 2:
				configObject.CDRepoChoice()
			case 3:
				configObject.CDRepoManual(cliArguments[2])
			default:
				fmt.Println("[Exit 1] Please select valid sub-action")
				os.Exit(1)
			}
			configObject.CloseFiles()
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
