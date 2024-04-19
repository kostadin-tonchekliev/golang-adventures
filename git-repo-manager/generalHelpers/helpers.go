package generalHelpers

import (
	"fmt"
	"git-repo-manager/sharedConstants"
	"github.com/cqroot/prompt"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	config2 "github.com/go-git/go-git/v5/config"
	"os"
)

var Yellow = color.New(color.FgYellow).SprintFunc()
var Blue = color.New(color.FgBlue).SprintFunc()
var Cyan = color.New(color.FgCyan).SprintFunc()
var Green = color.New(color.FgGreen).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()

// ReadInput - Read user input
func ReadInput(message string, prompt *prompt.Prompt, isFile bool) string {
	var (
		inputString string
		err         error
	)

	for {
		inputString, err = prompt.Ask(message).Input("")
		if err != nil {
			LogOutput(fmt.Sprintf("Unable to read input\n%s\n", err), 4, true)
		}

		if inputString != "" {
			if isFile {
				_, err = os.Stat(inputString)
				if err != nil {
					LogOutput("File/folder doesn't appear to exist", 4, false)
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
		LogOutput(fmt.Sprintf("Unable to read home directory\n%s\n", err), 4, true)
	}

	if _, err = os.Stat(fmt.Sprintf("%s/%s", homeDir, sharedConstants.ProjectHomeName)); err != nil {
		LogOutput(fmt.Sprintf("Project folder %s/%s cannot be found or isn't accessible. Has `setup` been ran?\n", homeDir, sharedConstants.ProjectHomeName), 4, true)
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
		"version": {
			Description: "Display build version of the executable",
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
			Description: "Add a single repository to the config",
			Level:       2,
			Origin:      "config",
		},
		"bulk-add": {
			Description: "Bulk add repositories to the config",
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

// DisplayVersion - Show the version of the executable
func DisplayVersion() {
	fmt.Printf("Build version: %s - %s\n", sharedConstants.BuildVersion, sharedConstants.BuildType)
	fmt.Printf("Build date: %s\n", sharedConstants.BuildDate)
}

// GetRepoUri - Get the uri of the remote repository
func GetRepoUri(repoPath string) string {
	var (
		repoObject     *git.Repository
		repoConfig     *config2.Config
		configMapValue *config2.RemoteConfig
		repoUri        string
		err            error
	)

	repoObject, err = git.PlainOpen(repoPath)
	if err != nil {
		LogOutput(fmt.Sprintf("[Err] Unable to initialize repository in %s\n%s\n", repoPath, err), 4, true)
	}

	repoConfig, err = repoObject.Config()
	if err != nil {
		LogOutput(fmt.Sprintf("[Err] Unable to detect config of repository %s\n%s\n", repoPath, err), 4, true)
	}

	for _, configMapValue = range repoConfig.Remotes {
		if configMapValue.Name == "origin" {
			repoUri = configMapValue.URLs[0]
		} else {
			repoUri = "Empty"
		}
	}

	return repoUri
}

// LogOutput - Display output
func LogOutput(outputText string, level int, exit bool) {
	// The existing levels are:
	// 1 - Debug
	// 2 - Info
	// 3 - Warning
	// 4 - Error
	switch level {
	case 1:
		fmt.Printf("[%s] %s\n", Blue("Debug"), outputText)
	case 2:
		fmt.Printf("[%s] %s\n", Cyan("Info"), outputText)
	case 3:
		fmt.Printf("[%s] %s\n", Yellow("Warn"), outputText)
	case 4:
		fmt.Printf("[%s] %s\n", Red("Err"), outputText)

	}

	if exit {
		if level == 4 {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
}
