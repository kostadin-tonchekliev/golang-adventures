package generalHelpers

import (
	"fmt"
	"git-repo-manager/sharedConstants"
	"github.com/cqroot/prompt"
	"os"
)

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

	homeDir = "exampleFiles" // Just for testing purposes

	if _, err = os.Stat(fmt.Sprintf("%s/%s", homeDir, sharedConstants.ProjectHomeName)); err != nil {
		fmt.Printf("[Err] Project folder %s/%s cannot be found or isn't accessible. Has `setup` been ran?\n%s\n", homeDir, sharedConstants.ProjectHomeName, err)
		os.Exit(1)
	}
}
