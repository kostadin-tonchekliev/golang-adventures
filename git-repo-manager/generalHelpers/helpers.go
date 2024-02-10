package generalHelpers

import (
	"fmt"
	"git-repo-manager/sharedConstants"
	"github.com/cqroot/prompt"
	"os"
)

// ReadInput - Read user input
func ReadInput(message string, prompt *prompt.Prompt, isFile bool) string {
	var (
		inputString string
		err         error
	)

	for {
		inputString, err = prompt.Ask(message).Input("")
		if err != nil {
			fmt.Println("[Err] Unable to read input\n", err)
			os.Exit(1)
		}

		if inputString != "" {
			if isFile {
				_, err = os.Stat(inputString)
				if err != nil {
					fmt.Println("[Err] File/folder doesn't appear to exist")
				} else {
					break
				}
			} else {
				break
			}
		}
	}
	return inputString
}

// VerifyEnv - Verify that the project environment is properly setup
func VerifyEnv() {
	var (
		homeDir string
		err     error
	)

	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println("[Err] Unable to read home directory\n", err)
		os.Exit(1)
	}

	if _, err = os.Stat(fmt.Sprintf("%s/%s", homeDir, sharedConstants.ProjectHomeName)); err != nil {
		fmt.Printf("[Err] Project folder %s/%s cannot be found or isn't accessible. Has `setup` been ran?\n%s\n", homeDir, sharedConstants.ProjectHomeName, err)
		os.Exit(1)
	}
}

// ShowHelp - Display the supported arguments
func ShowHelp(level int, origin string) {
	type arg struct {
		Description string
		Level       int
		Origin      string
	}

	arguments := map[string]arg{
		"setup": {
			Description: "Configure project environment",
			Level:       1,
			Origin:      "",
		},
		"status": {
			Description: "Show status of configured repositories",
			Level:       1,
			Origin:      "",
		},
		"config": {
			Description: "Actions related to the config file",
			Level:       1,
			Origin:      "",
		},
		"ls": {
			Description: "List existing config",
			Level:       2,
			Origin:      "config",
		},
		"add": {
			Description: "Add new repositories to config",
			Level:       2,
			Origin:      "config",
		},
		"remove": {
			Description: "Remove repositories from config",
			Level:       2,
			Origin:      "config",
		},
		"cd": {
			Description: "Move to a repository. Not to be used directly!",
			Level:       1,
			Origin:      "",
		},
		"empty": {
			Description: "Choose the desired repository from a list",
			Level:       2,
			Origin:      "cd",
		},
		"repository_petname": {
			Description: "Go to the repository which matches the petname",
			Level:       2,
			Origin:      "cd",
		},
	}

	fmt.Println("Usage:")
	for argument, content := range arguments {
		if level == content.Level && origin == content.Origin {
			fmt.Printf("\t %s - %s\n", argument, content.Description)
		}
	}
}
