package generalHelpers

import (
	"fmt"
	"git-repo-manager/sharedConstants"
	"github.com/cqroot/prompt"
	"os"
)

func ReadInput(message string, prompt *prompt.Prompt) string {
	var inputString string
	var err error

	for {
		inputString, err = prompt.Ask(message).Input("")
		if err != nil {
			fmt.Println("[Err] Unable to read input\n", err)
			os.Exit(0)
		}

		if inputString != "" {
			break
		}
	}
	return inputString
}

func SetupEnv() {
	var homeDir string
	var err error
	var fileObject *os.File

	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println("[Err] Unable to read home directory\n", err)
		os.Exit(1)
	}

	homeDir = "exampleFiles" // Just for testing purposes

	if _, err = os.Stat(fmt.Sprintf("%s/%s", homeDir, sharedConstants.ProjectHomeName)); err != nil {
		err := os.Mkdir(fmt.Sprintf("%s/%s", homeDir, sharedConstants.ProjectHomeName), 0755)
		if err != nil {
			fmt.Printf("[Err] Unable to create project folder %s/%s \n%s\n", homeDir, sharedConstants.ProjectHomeName, err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("[Info] Project folder %s/%s already exists\n", homeDir, sharedConstants.ProjectHomeName)
	}

	fmt.Printf("[Info] Project home folder initialized in %s/%s \n", homeDir, sharedConstants.ProjectHomeName)

	for _, fileName := range []string{sharedConstants.ConfigFileName, sharedConstants.TmpDirFileName} {
		if _, err = os.Stat(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, fileName)); err != nil {
			fileObject, err = os.Create(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, fileName))
			if err != nil {
				fmt.Printf("[Err] Unable to create file %s/%s/%s \n%s\n", homeDir, sharedConstants.ProjectHomeName, fileName, err)
				os.Exit(0)
			}

			fmt.Printf("[Info] File %s/%s/%s succesfully created\n", homeDir, sharedConstants.ProjectHomeName, fileName)

			defer fileObject.Close()
			err = fileObject.Chmod(0644)
			if err != nil {
				fmt.Printf("[Err] Unable to change permissions of file  %s/%s/%s \n%s\n", homeDir, sharedConstants.ProjectHomeName, fileName, err)
			}
		} else {
			fmt.Printf("[Info] File %s/%s/%s already exists\n", homeDir, sharedConstants.ProjectHomeName, fileName)
		}
	}
}

func VerifyEnv() {
	var homeDir string
	var err error

	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println("[Err] Unable to read home directory\n", err)
		os.Exit(1)
	}

	homeDir = "exampleFiles" // Just for testing purposes

	if _, err = os.Stat(fmt.Sprintf("%s/%s", homeDir, sharedConstants.ProjectHomeName)); err != nil {
		fmt.Printf("[Err] Project folder %s/%s cannot be found or isn't accessible. Has `setup` been ran?\n%s\n", homeDir, sharedConstants.ProjectHomeName, err)
		os.Exit(0)
	}
}
