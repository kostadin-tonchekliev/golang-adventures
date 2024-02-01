package configActions

import (
	"encoding/json"
	"fmt"
	"git-repo-manager/generalHelpers"
	"git-repo-manager/sharedConstants"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/multichoose"
	"os"
)

var promptObject = prompt.New() // Can most likely remove this if I don't end up using it again

type Config struct {
	RepoMap    map[string]RepoObject
	ConfigFile *os.File
	TmpDirFile *os.File
}

type RepoObject struct {
	Url         string `json:"url"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

func SetupEnv() {
	var (
		homeDir, fileName, execLocation, shellType, shellFileLocation string
		fileObject                                                    *os.File
		err                                                           error
	)

	homeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Println("[Err] Unable to read home directory\n", err)
		os.Exit(1)
	}

	homeDir = "exampleFiles" // Just for testing purposes

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

	execLocation = generalHelpers.ReadInput("Select path of executable:", promptObject, true)
	shellType = os.Getenv("SHELL")

	switch shellType {
	case "/bin/zsh":
		shellFileLocation = fmt.Sprintf("%s/.zshrc", homeDir)
	case "/bin/bash":
		shellFileLocation = fmt.Sprintf("%s/.bashrc", homeDir)
	default:
		fmt.Println("[Err] Unknown shell type:", shellType)
		os.Exit(1)
	}

	fileObject, err = os.OpenFile(shellFileLocation, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("[Err] Unable to open file\n", err)
		os.Exit(1)
	}
	// _, err = fileObject.WriteString(fmt.Sprintf("\nfunction %s() {%s cd $1; if [[ $? == 0 ]]; then cd $(cat %s/%s/%s); fi}\n", sharedConstants.AliasName, execLocation, homeDir, sharedConstants.ProjectHomeName, sharedConstants.TmpDirFileName))
	_, err = fileObject.WriteString(fmt.Sprintf("\nfunction %s() {go run %s cd $1; if [[ $? == 0 ]]; then cd $(cat %s/%s/%s); fi}\n", sharedConstants.AliasName, execLocation, homeDir, sharedConstants.ProjectHomeName, sharedConstants.TmpDirFileName)) // Just for testing purposes
	if err != nil {
		fmt.Println("[Err] Unable to write to file\n", err)
	}

	fmt.Printf("[Info] Please run `source %s` or reload your terminal\n", shellFileLocation)
}

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

	homeDir = "exampleFiles" // Just for testing purposes

	fmt.Printf("[Info] Reading config file located in %s/%s\n", homeDir, sharedConstants.ConfigFileName)
	configFileObject, err = os.OpenFile(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, sharedConstants.ConfigFileName), os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("[Err] Unable to read/create config file\n", err)
		os.Exit(1)
	}

	fmt.Printf("[Info] Reading temporary directory file located in %s/%s\n", homeDir, sharedConstants.TmpDirFileName)
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

func (config Config) CloseFiles() {
	var (
		err error
	)

	fmt.Println("[Info] Closing config file", config.ConfigFile.Name())
	err = config.ConfigFile.Close()
	if err != nil {
		fmt.Println("[Err] Unable to close config file\n", err)
		os.Exit(1)
	}

	fmt.Println("[Info] Closing temporary directory file", config.TmpDirFile.Name())
	err = config.TmpDirFile.Close()
	if err != nil {
		fmt.Println("[Err] Unable to close temporary directory file\n", err)
		os.Exit(1)
	}
}

func (config Config) ListConfig() {
	var (
		repoContent RepoObject
		petName     string
		counter     int = 1
	)

	if len(config.RepoMap) != 0 {
		for petName, repoContent = range config.RepoMap {
			fmt.Printf("[%d] %s\n", counter, petName)
			fmt.Println("\tURL:", repoContent.Url)
			fmt.Println("\tPath:", repoContent.Path)
			fmt.Println("\tDescription:", repoContent.Description)
			fmt.Println()
			counter += 1
		}
	} else {
		fmt.Println("[Warn] Config file empty")
		os.Exit(0)
	}
}

func (config Config) AddConfig() {
	var (
		petName, uri, path, description string
		jsonContent                     []byte
		err                             error
	)

	petName = generalHelpers.ReadInput("Select petname for the repository", promptObject, false)
	description = generalHelpers.ReadInput("Select description for the repository", promptObject, false)
	uri = generalHelpers.ReadInput("Select uri of the repository", promptObject, false)
	path = generalHelpers.ReadInput("Select local path of the repository", promptObject, true)

	config.RepoMap[petName] = RepoObject{
		Url:         uri,
		Path:        path,
		Description: description,
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

func (config Config) RemoveConfig() {
	var (
		choiceOptions, repoSelection []string
		petName, repository          string
		jsonContent                  []byte
		err                          error
	)

	for petName, _ = range config.RepoMap {
		choiceOptions = append(choiceOptions, petName)
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

func (config Config) CDRepoChoice() {
	var (
		choiceOptions          []choose.Choice
		repoObject             choose.Choice
		petName, repoSelection string
		repoContent            RepoObject
		err                    error
	)

	for petName, repoContent = range config.RepoMap {
		repoObject = choose.Choice{
			Text: petName,
			Note: repoContent.Url,
		}

		choiceOptions = append(choiceOptions, repoObject)
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
