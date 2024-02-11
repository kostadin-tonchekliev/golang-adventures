package main

import (
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
		case "status":
			generalHelpers.VerifyEnv()
			configObject := configActions.ReadConfig()
			configObject.RepoStatus()
			configObject.CloseFiles()
		case "version":
			generalHelpers.DisplayVersion()
		case "config":
			generalHelpers.VerifyEnv()
			configObject := configActions.ReadConfig()
			if len(cliArguments) > 2 {
				switch cliArguments[2] {
				case "ls":
					configObject.ListConfig()
				case "add":
					configObject.AddConfig()
				case "remove":
					configObject.RemoveConfig()
				default:
					generalHelpers.ShowHelp(2, "config")
					os.Exit(1)
				}
			} else {
				generalHelpers.ShowHelp(2, "config")
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
				generalHelpers.ShowHelp(2, "cd")
				os.Exit(1)
			}
			configObject.CloseFiles()
		default:
			generalHelpers.ShowHelp(1, "")
			os.Exit(1)
		}
	} else {
		generalHelpers.ShowHelp(1, "")
		os.Exit(1)
	}
}
