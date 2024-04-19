package configActions

import (
	"encoding/json"
	"fmt"
	"git-repo-manager/generalHelpers"
	"git-repo-manager/sharedConstants"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/multichoose"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"io/fs"
	"os"
	"slices"
)

// Objects used through the script
var promptObject = prompt.New()

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
		generalHelpers.LogOutput("Unable to read the home directory", 3, true)
	}

	execLocation, err = os.Executable()
	if err != nil {
		generalHelpers.LogOutput("Unable to find the path of the executable, please enter it manually.", 3, false)
		execLocation = generalHelpers.ReadInput("Select the path of the executable. It needs to be absolute path!", promptObject, true)
	}

	if _, err = os.Stat(fmt.Sprintf("%s/%s", homeDir, sharedConstants.ProjectHomeName)); err != nil {
		err = os.Mkdir(fmt.Sprintf("%s/%s", homeDir, sharedConstants.ProjectHomeName), 0755)
		if err != nil {
			generalHelpers.LogOutput(fmt.Sprintf("Unable to create project folder %s/%s \n%s\n", homeDir, sharedConstants.ProjectHomeName, err), 3, true)
		}
		generalHelpers.LogOutput(fmt.Sprintf("Project home folder initialized in %s/%s \n", homeDir, sharedConstants.ProjectHomeName), 2, false)
	} else {
		generalHelpers.LogOutput(fmt.Sprintf("The project already appears to be setup. Existing config found in %s/%s\n", homeDir, sharedConstants.ProjectHomeName), 2, false)
	}

	for _, fileName = range []string{sharedConstants.ConfigFileName, sharedConstants.TmpDirFileName} {
		if _, err = os.Stat(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, fileName)); err != nil {
			fileObject, err = os.Create(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, fileName))
			if err != nil {
				generalHelpers.LogOutput(fmt.Sprintf("Unable to create file %s/%s/%s \n%s\n", homeDir, sharedConstants.ProjectHomeName, fileName, err), 4, true)
			}

			generalHelpers.LogOutput(fmt.Sprintf("The file %s/%s/%s is succesfully created\n", homeDir, sharedConstants.ProjectHomeName, fileName), 2, false)

			defer fileObject.Close()
			err = fileObject.Chmod(0644)
			if err != nil {
				generalHelpers.LogOutput(fmt.Sprintf("Unable to change permissions of file %s/%s/%s \n%s\n", homeDir, sharedConstants.ProjectHomeName, fileName, err), 4, true)
			}
		} else {
			generalHelpers.LogOutput(fmt.Sprintf("The file %s/%s/%s already exists\n", homeDir, sharedConstants.ProjectHomeName, fileName), 2, false)
		}
	}

	homeContent, err = os.ReadDir(homeDir)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to read the content of the home directory in %s\n%s\n", homeDir, err), 4, true)
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
		generalHelpers.LogOutput(fmt.Sprintf("Unable to select choice\n%s\n", err), 4, true)
	}

	functionContent = fmt.Sprintf("\nfunction %s() { %s cd $1; if [[ $? == 0 ]]; then cd $(cat %s/%s/%s); fi }\n", sharedConstants.AliasName, execLocation, homeDir, sharedConstants.ProjectHomeName, sharedConstants.TmpDirFileName)

	switch configChoice {
	case "Custom":
		configLocation = generalHelpers.ReadInput("Select your desired config. It needs an absolute path!", promptObject, true)
	case "Skip":
		generalHelpers.LogOutput(fmt.Sprintf("Append the following function to your desired file, this will allow you to cd directly into repositories:\n%s\n", functionContent), 2, true)
	default:
		configLocation = fmt.Sprintf("%s/%s", homeDir, configChoice)
	}

	fileObject, err = os.OpenFile(configLocation, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to open the following file for writing: %s\n%s", configLocation, err), 4, true)
	}

	defer fileObject.Close()
	_, err = fileObject.WriteString(functionContent)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to write function to config:\n%s\n", err), 4, true)
	}

	generalHelpers.LogOutput("Write successfull!", 2, false)
	generalHelpers.LogOutput(fmt.Sprintf("Please run `source %s` or reload your terminal\n", configLocation), 2, true)
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
		generalHelpers.LogOutput(fmt.Sprintf("Unable to get the home directory\n%s\n", err), 4, true)

	}

	configFileObject, err = os.OpenFile(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, sharedConstants.ConfigFileName), os.O_RDWR, 0644)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to read/create config file\n%s\n", err), 4, true)
	}

	tmpDirFileObject, err = os.OpenFile(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, sharedConstants.TmpDirFileName), os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to read/create temporary directory file\n%s\n", err), 4, true)
	}

	configFileStat, _ = configFileObject.Stat()
	if configFileStat.Size() == 0 {
		configFileObject.Write([]byte(`{}`))
		configFileObject.Close()
		configFileObject, err = os.OpenFile(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, sharedConstants.ConfigFileName), os.O_RDWR, 0644)
		if err != nil {
			generalHelpers.LogOutput(fmt.Sprintf("Unable to read/create config file\n%s\n", err), 4, true)
		}
	}

	configFileStat, _ = configFileObject.Stat()
	buffer = make([]byte, configFileStat.Size())
	_, err = configFileObject.Read(buffer)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to read config file\n%s\n", err), 4, true)

	}

	err = json.Unmarshal(buffer, &configFile.RepoMap)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to unmarshal json\n%s\n", err), 4, true)
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
		generalHelpers.LogOutput(fmt.Sprintf("Unable to close config file\n%s\n", err), 4, true)
	}

	err = config.TmpDirFile.Close()
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to close temporary directory file\n%s\n", err), 4, true)
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
				generalHelpers.LogOutput(fmt.Sprintf("Unable to open %s located in %s\n%s", petName, repoContent.Path, err), 4, true)

			}

			branchName, err = repoObject.Head()
			if err != nil {
				generalHelpers.LogOutput(fmt.Sprintf("Unable to read branch\n%s\n", err), 4, true)
			}

			repoWorktree, err = repoObject.Worktree()
			if err != nil {
				generalHelpers.LogOutput(fmt.Sprintf("Unable to read worktree\n%s\n", err), 4, true)
			}

			repoStatus, err = repoWorktree.Status()
			if err != nil {
				generalHelpers.LogOutput(fmt.Sprintf("Unable to read repo status\n%s\n", err), 4, true)
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

			fmt.Printf("[%s] in ", generalHelpers.Blue(petName))
			if branchName.Name() == "refs/heads/main" || branchName.Name() == "refs/heads/master" {
				fmt.Printf("%s:\n", generalHelpers.Yellow(branchName.Name()))
			} else {
				fmt.Printf("%s:\n", generalHelpers.Green(branchName.Name()))
			}

			for fileType, fileCount := range fileMap {
				fmt.Printf("\t%s: %d\n", fileType, fileCount)
			}
		}
	} else {
		generalHelpers.LogOutput("Config file empty", 3, true)
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
			fmt.Printf("[%s] in %s: %s\n", generalHelpers.Blue(petName), generalHelpers.Yellow(repoContent.Path), generalHelpers.Cyan(repoContent.Url))
		}
	} else {
		generalHelpers.LogOutput("Config file empty", 3, true)
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
		generalHelpers.LogOutput(fmt.Sprintf("Unable to read the content of %s\n%s\n", repoPath, err), 4, true)
	}

	for _, dirEntry = range dirContent {
		if dirEntry.IsDir() {
			singleRepoPath = fmt.Sprintf("%s/%s", repoPath, dirEntry.Name())
			generalHelpers.LogOutput(fmt.Sprintf("Found subdirectory: %s. Verifying if it contains a valid repository...", dirEntry.Name()), 2, false)
			repoObject, err = git.PlainOpen(singleRepoPath)
			if err != nil {
				generalHelpers.LogOutput(fmt.Sprintf("%s doesn't appear to contain a valid repository\n", dirEntry.Name()), 3, false)
			} else {
				// Just a quick command to verify if it is working
				_, err = repoObject.Head()
				if err != nil {
					generalHelpers.LogOutput(fmt.Sprintf("The subdirectory %s appears to contain a repository but coun't get HEAD. Still going to add the repository, but manually verify if it should be added\n%s\n", dirEntry.Name(), err), 3, false)
				}
				config.AddConfig(dirEntry.Name(), singleRepoPath)
				generalHelpers.LogOutput(fmt.Sprintf("Successfully added %s to the config\n", dirEntry.Name()), 2, false)
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
		generalHelpers.LogOutput(fmt.Sprintf("Unable to convert map to json\n%s\n", err), 4, true)
	}

	config.ConfigFile.Truncate(0)
	config.ConfigFile.Seek(0, 0)
	_, err = config.ConfigFile.Write(jsonContent)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to write new value to config\n%s\n", err), 4, true)
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
		generalHelpers.LogOutput("Config file empty", 3, true)
	}

	repoSelection, err = promptObject.Ask("Select repositories to remove").MultiChoose(choiceOptions, multichoose.WithHelp(true))
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Error reading repository choice\n%s\n", err), 4, true)
	}

	for _, repository = range repoSelection {
		delete(config.RepoMap, repository)
	}

	jsonContent, err = json.MarshalIndent(config.RepoMap, "", "  ")
	config.ConfigFile.Truncate(0)
	config.ConfigFile.Seek(0, 0)
	_, err = config.ConfigFile.Write(jsonContent)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to write new value to config\n%s\n", err), 4, true)
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
		generalHelpers.LogOutput("Config file empty", 3, true)
	}

	repoSelection, err = promptObject.Ask("Select repository").AdvancedChoose(choiceOptions)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Error reading repository choice\n%s\n", err), 4, true)
	}

	_, err = config.TmpDirFile.WriteString(config.RepoMap[repoSelection].Path)
	if err != nil {
		generalHelpers.LogOutput(fmt.Sprintf("Unable to write to temporary directory file\n%s\n", err), 4, true)
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
			generalHelpers.LogOutput(fmt.Sprintf("Unable to write to temporary directory file\n%s\n", err), 4, true)
		}
	} else {
		generalHelpers.LogOutput(fmt.Sprintf("Repository %s not found in the config\n", repoSelection), 4, true)
	}
}
