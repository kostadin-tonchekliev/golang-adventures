package configActions

import (
	"encoding/json"
	"fmt"
	"git-repo-manager/generalHelpers"
	"git-repo-manager/sharedConstants"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/multichoose"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"io/fs"
	"os"
	"slices"
)

// Objects used through the script
var promptObject = prompt.New()
var yellow = color.New(color.FgYellow).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()

type Config struct {
	RepoMap    map[string]RepoObject
	ConfigFile *os.File
	TmpDirFile *os.File
}

type RepoObject struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

// SetupEnv - Build and verify project environment
func SetupEnv() {
	var (
		execLocation, homeDir, fileName, configChoice, configLocation, functionContent string
		homeContent                                                                    []os.DirEntry
		singleFile                                                                     fs.DirEntry
		fileObject                                                                     *os.File
		supportedConfigFiles, existingConfigFiles                                      []string
		err                                                                            error
	)

	// Define supported config files
	supportedConfigFiles = []string{".bashrc", ".bash_alias", ".bash_functions", ".zshrc", ".zsh_alias", ".zsh_functions"}

	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println("[Err] Unable to read home directory\n", err)
		os.Exit(1)
	}

	execLocation, err = os.Executable()
	if err != nil {
		fmt.Println("[Err] Unable to find path of executable, please enter it manually")
		execLocation = generalHelpers.ReadInput("Select path of executable:", promptObject, true)
	}

	if _, err = os.Stat(fmt.Sprintf("%s/%s", homeDir, sharedConstants.ProjectHomeName)); err != nil {
		err = os.Mkdir(fmt.Sprintf("%s/%s", homeDir, sharedConstants.ProjectHomeName), 0755)
		if err != nil {
			fmt.Printf("[Err] Unable to create project folder %s/%s \n%s\n", homeDir, sharedConstants.ProjectHomeName, err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("[Info] Project folder %s/%s already exists\n", homeDir, sharedConstants.ProjectHomeName)
	}

	fmt.Printf("[Info] Project home folder initialized in %s/%s \n", homeDir, sharedConstants.ProjectHomeName)

	for _, fileName = range []string{sharedConstants.ConfigFileName, sharedConstants.TmpDirFileName} {
		if _, err = os.Stat(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, fileName)); err != nil {
			fileObject, err = os.Create(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, fileName))
			if err != nil {
				fmt.Printf("[Err] Unable to create file %s/%s/%s \n%s\n", homeDir, sharedConstants.ProjectHomeName, fileName, err)
				os.Exit(1)
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

	homeContent, err = os.ReadDir(homeDir)
	if err != nil {
		fmt.Printf("Unable to read the content of %s\n%s\n", homeDir, err)
		os.Exit(0)
	}

	for _, singleFile = range homeContent {
		if !singleFile.IsDir() {
			if slices.Contains(supportedConfigFiles, singleFile.Name()) {
				existingConfigFiles = append(existingConfigFiles, singleFile.Name())
			}
		}
	}

	// Define custom options
	existingConfigFiles = append(existingConfigFiles, "Custom", "Skip")

	configChoice, err = promptObject.Ask("Select config file for cd function:").Choose(existingConfigFiles)
	if err != nil {
		fmt.Println("[Err] Unable to select choice\n", err)
	}

	functionContent = fmt.Sprintf("\nfunction %s() { %s cd $1; if [[ $? == 0 ]]; then cd $(cat %s/%s/%s); fi }\n", sharedConstants.AliasName, execLocation, homeDir, sharedConstants.ProjectHomeName, sharedConstants.TmpDirFileName)

	switch configChoice {
	case "Custom":
		configLocation = generalHelpers.ReadInput("Select your desired config", promptObject, true)
	case "Skip":
		fmt.Println("Append the following function to your desired file, this will allow you to cd directly into repositories:\n", functionContent)
		os.Exit(1)
	default:
		configLocation = configChoice
	}

	fileObject, err = os.OpenFile(fmt.Sprintf("%s/%s", homeDir, configLocation), os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("[Err] Unable to open config file for writing:\n", err)
	}

	defer fileObject.Close()
	_, err = fileObject.WriteString(functionContent)
	if err != nil {
		fmt.Println("[Err] Unable to write function to config:\n", err)
	}

	fmt.Println("[Info] Write successfull!")
	fmt.Printf("[Info] Please run `source %s` or reload your terminal\n", configLocation)
}

// ReadConfig - Read the config file and build it into the Config struct
func ReadConfig() Config {
	var (
		configFileObject, tmpDirFileObject *os.File
		configFileStat                     os.FileInfo
		configFile                         Config
		homeDir                            string
		buffer                             []byte
		err                                error
	)

	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println("[Err] Unable to read home directory\n", err)
		os.Exit(1)
	}

	configFileObject, err = os.OpenFile(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, sharedConstants.ConfigFileName), os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("[Err] Unable to read/create config file\n", err)
		os.Exit(1)
	}

	tmpDirFileObject, err = os.OpenFile(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, sharedConstants.TmpDirFileName), os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("[Err] Unable to read/create temporary directory file\n", err)
		os.Exit(1)
	}

	configFileStat, _ = configFileObject.Stat()
	if configFileStat.Size() == 0 {
		configFileObject.Write([]byte(`{}`))
		configFileObject.Close()
		configFileObject, err = os.OpenFile(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, sharedConstants.ConfigFileName), os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("[Err] Unable to read/create config file\n", err)
			os.Exit(1)
		}
	}

	configFileStat, _ = configFileObject.Stat()
	buffer = make([]byte, configFileStat.Size())
	_, err = configFileObject.Read(buffer)
	if err != nil {
		fmt.Println("[Err] Unable to read config file\n", err)
		os.Exit(1)
	}

	err = json.Unmarshal(buffer, &configFile.RepoMap)
	if err != nil {
		fmt.Println("[Err] Unable to unmarshal json\n", err)
		os.Exit(1)
	}

	configFile.ConfigFile = configFileObject
	configFile.TmpDirFile = tmpDirFileObject

	return configFile
}

// CloseFiles - Close all files opened by other methods
func (config Config) CloseFiles() {
	var (
		err error
	)

	err = config.ConfigFile.Close()
	if err != nil {
		fmt.Println("[Err] Unable to close config file\n", err)
		os.Exit(1)
	}

	err = config.TmpDirFile.Close()
	if err != nil {
		fmt.Println("[Err] Unable to close temporary directory file\n", err)
		os.Exit(1)
	}
}

// RepoStatus - Print the status of the repositories based on the local changes
func (config Config) RepoStatus() {
	var (
		repoContent  RepoObject
		repoObject   *git.Repository
		repoWorktree *git.Worktree
		repoStatus   git.Status
		fileStatus   *git.FileStatus
		branchName   *plumbing.Reference
		petName      string
		err          error
	)

	if len(config.RepoMap) != 0 {
		for petName, repoContent = range config.RepoMap {
			fileMap := make(map[string]int)

			repoObject, err = git.PlainOpen(repoContent.Path)
			if err != nil {
				fmt.Println("[Err] Unable to open repository\n", err)
				os.Exit(1)
			}

			branchName, err = repoObject.Head()
			if err != nil {
				fmt.Println("[Err] Unable to read branch\n", err)
				os.Exit(1)
			}

			repoWorktree, err = repoObject.Worktree()
			if err != nil {
				fmt.Println("[Err] Unable to read worktree\n", err)
			}

			repoStatus, err = repoWorktree.Status()
			if err != nil {
				fmt.Println("[Err] Unable to read repo status\n", err)
			}

			for _, fileStatus = range repoStatus {
				switch string(fileStatus.Staging) {
				case "?":
					fileMap["Untracked"]++
				case "M":
					fileMap["Modified"]++
				case "A":
					fileMap["Added"]++
				case "D":
					fileMap["Deleted"]++
				case "R":
					fileMap["Renamed"]++
				case "C":
					fileMap["Copied"]++
				}
			}

			fmt.Printf("[%s] in ", blue(petName))
			if branchName.Name() == "refs/heads/main" || branchName.Name() == "refs/heads/master" {
				fmt.Printf("%s:\n", yellow(branchName.Name()))
			} else {
				fmt.Printf("%s:\n", green(branchName.Name()))
			}

			for fileType, fileCount := range fileMap {
				fmt.Printf("\t%s: %d\n", fileType, fileCount)
			}
		}
	} else {
		fmt.Println("[Warn] Config file empty")
		os.Exit(0)
	}
}

// ListConfig - Display the content of the config file
func (config Config) ListConfig() {
	var (
		repoContent RepoObject
		petName     string
	)

	if len(config.RepoMap) != 0 {
		for petName, repoContent = range config.RepoMap {
			fmt.Printf("[%s] in %s: %s\n", blue(petName), yellow(repoContent.Path), cyan(repoContent.Url))
		}
	} else {
		fmt.Println("[Warn] Config file empty")
		os.Exit(0)
	}
}

// AddSingleConfig - Add additional entries to the config file
func (config Config) AddSingleConfig() {
	var (
		petName, path string
	)
	petName = generalHelpers.ReadInput("Select petname for the repository", promptObject, false)
	path = generalHelpers.ReadInput("Select absolute path of the repository", promptObject, true)

	config.AddConfig(petName, path)
}

// AddMultipleConfig - Add automatically multiple repositories by specifying a single main directory
func (config Config) AddMultipleConfig() {
	var (
		repoPath, singleRepoPath string
		repoObject               *git.Repository
		dirContent               []os.DirEntry
		dirEntry                 os.DirEntry
		err                      error
	)

	repoPath = generalHelpers.ReadInput("Select the absolute path where all repositories are stored", promptObject, true)

	dirContent, err = os.ReadDir(repoPath)
	if err != nil {
		fmt.Printf("[Err] Unable to read the content of %s\n%s\n", repoPath, err)
		os.Exit(1)
	}

	for _, dirEntry = range dirContent {
		if dirEntry.IsDir() {
			singleRepoPath = fmt.Sprintf("%s/%s", repoPath, dirEntry.Name())
			fmt.Printf("[Info] Found subdirectory: %s. Verifying if it contains a valid repository\n", dirEntry.Name())
			repoObject, err = git.PlainOpen(singleRepoPath)
			if err != nil {
				fmt.Printf("[Warn] %s doesn't appear to contain a valid repository\n", dirEntry.Name())
			} else {
				// Just a quick command to verify if it is working
				_, err = repoObject.Head()
				if err != nil {
					fmt.Printf("[Warn] The subdirectory %s appears to contain a repository but coun't get HEAD\n%s\n", dirEntry.Name(), err)
				}
				config.AddConfig(dirEntry.Name(), singleRepoPath)
				fmt.Printf("[Info] Successfully added %s to the config\n", dirEntry.Name())
			}
		}
	}
}

// AddConfig - Add config into the config file
func (config Config) AddConfig(petName string, path string) {
	var (
		jsonContent []byte
		uri         string
		err         error
	)

	uri = generalHelpers.GetRepoUri(path)

	config.RepoMap[petName] = RepoObject{
		Url:  uri,
		Path: path,
	}

	jsonContent, err = json.MarshalIndent(config.RepoMap, "", "  ")
	if err != nil {
		fmt.Println("[Err] Unable to convert map to json\n", err)
	}

	config.ConfigFile.Truncate(0)
	config.ConfigFile.Seek(0, 0)
	_, err = config.ConfigFile.Write(jsonContent)
	if err != nil {
		fmt.Println("[Err] Unable to write new value to config\n", err)
	}
}

// RemoveConfig - Remove entries from the config file
func (config Config) RemoveConfig() {
	var (
		choiceOptions, repoSelection []string
		petName, repository          string
		jsonContent                  []byte
		err                          error
	)

	if len(config.RepoMap) > 0 {
		for petName = range config.RepoMap {
			choiceOptions = append(choiceOptions, petName)
		}
	} else {
		fmt.Println("[Warn] Config file empty")
		os.Exit(0)
	}

	repoSelection, err = promptObject.Ask("Select repositories to remove").MultiChoose(choiceOptions, multichoose.WithHelp(true))
	if err != nil {
		fmt.Println("[Err] Error reading repository choice\n", err)
		os.Exit(1)
	}

	for _, repository = range repoSelection {
		delete(config.RepoMap, repository)
	}

	jsonContent, err = json.MarshalIndent(config.RepoMap, "", "  ")
	config.ConfigFile.Truncate(0)
	config.ConfigFile.Seek(0, 0)
	_, err = config.ConfigFile.Write(jsonContent)
	if err != nil {
		fmt.Println("[Err] Unable to write new value to config\n", err)
	}
}

// CDRepoChoice - Move to a repository which is selected from a list
func (config Config) CDRepoChoice() {
	var (
		choiceOptions          []choose.Choice
		repoObject             choose.Choice
		petName, repoSelection string
		repoContent            RepoObject
		err                    error
	)

	if len(config.RepoMap) > 0 {
		for petName, repoContent = range config.RepoMap {
			repoObject = choose.Choice{
				Text: petName,
				Note: repoContent.Url,
			}

			choiceOptions = append(choiceOptions, repoObject)
		}
	} else {
		fmt.Println("[Warn] Config file empty")
		os.Exit(1)
	}

	repoSelection, err = promptObject.Ask("Select repository").AdvancedChoose(choiceOptions)
	if err != nil {
		fmt.Println("[Err] Error reading repository choice\n", err)
		os.Exit(1)
	}

	_, err = config.TmpDirFile.WriteString(config.RepoMap[repoSelection].Path)
	if err != nil {
		fmt.Println("[Err] Unable to write to temporary directory file\n", err)
	}
}

// CDRepoManual -  Move to repository which has been manually selected
func (config Config) CDRepoManual(repoSelection string) {
	var (
		exist bool
		err   error
	)

	_, exist = config.RepoMap[repoSelection]
	if exist {
		_, err = config.TmpDirFile.WriteString(config.RepoMap[repoSelection].Path)
		if err != nil {
			fmt.Println("[Err] Unable to write to temporary directory file\n", err)
		}
	} else {
		fmt.Printf("[Err] Repository %s not found in the config\n", repoSelection)
		os.Exit(1)
	}
}
