package configActions

import (
	"encoding/json"
	"fmt"
	"git-repo-manager/generalHelpers"
	"git-repo-manager/sharedConstants"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"os"
)

var promptObject = prompt.New() // Can most likely remove this if I don't end up using it again

type Config struct {
	RepoMap    map[string]RepoObject
	ConfigFile *os.File
}

type RepoObject struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

func ReadConfig() Config {
	var configFile Config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("[Err] Unable to read home directory\n", err)
		os.Exit(1)
	}

	homeDir = "exampleFiles" // Just for testing purposes

	// fmt.Printf("[Info] Reading config file located in %s/%s\n", homeDir, configFileName) // Enable only for logging purposes

	configFileObject, err := os.OpenFile(fmt.Sprintf("%s/%s/%s", homeDir, sharedConstants.ProjectHomeName, sharedConstants.ConfigFileName), os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("[Err] Unable to read/create config file\n", err)
		os.Exit(1)
	}

	configFileStat, _ := configFileObject.Stat()
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
	buffer := make([]byte, configFileStat.Size())
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

	return configFile
}

func (config Config) CloseConfig() {
	// fmt.Println("[Info] Closing config file", config.ConfigFile.Name()) // Enable only for logging purposes
	err := config.ConfigFile.Close()
	if err != nil {
		fmt.Println("[Err] Unable to close file\n", err)
		os.Exit(1)
	}
}

func (config Config) ListConfig() {
	// fmt.Println("[Info] Reading content of config", config.ConfigFile.Name()) // Enable only for logging purposes
	i := 1
	if len(config.RepoMap) != 0 {
		for petName, repoContent := range config.RepoMap {
			fmt.Printf("[%d] %s\n", i, petName)
			fmt.Println("\tURL:", repoContent.Url)
			fmt.Println("\tPath:", repoContent.Path)
			fmt.Println()
			i += 1
		}
	} else {
		fmt.Println("[Warn] Config file empty")
		os.Exit(0)
	}
}

func (config Config) AddConfig() {
	var petName, uri, path string

	petName = generalHelpers.ReadInput("Select petname for the repository", promptObject)
	uri = generalHelpers.ReadInput("Select uri of the repository", promptObject)

	for {
		path = generalHelpers.ReadInput("Select local path of the repository", promptObject)
		files, err := os.ReadDir(path)
		if err != nil || len(files) == 0 {
			fmt.Println("[Err] Please select valid directory")
		} else {
			break
		}
	}

	config.RepoMap[petName] = RepoObject{
		Url:  uri,
		Path: path,
	}

	jsonContent, err := json.MarshalIndent(config.RepoMap, "", "  ")
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

func (config Config) CDRepoChoice() {
	var choiceOptions []choose.Choice

	for petName, repoContent := range config.RepoMap {
		repoChoice := choose.Choice{
			Text: petName,
			Note: repoContent.Url,
		}

		choiceOptions = append(choiceOptions, repoChoice)
	}

	repoChoice, err := promptObject.Ask("Select repository").AdvancedChoose(choiceOptions)
	if err != nil {
		fmt.Println("[Err] Error reading repository choice\n", err)
		os.Exit(1)
	}

	fmt.Print(config.RepoMap[repoChoice].Path)
}

func (config Config) CDRepoManual(repoChoice string) {
	_, exist := config.RepoMap[repoChoice]
	if exist {
		fmt.Print(config.RepoMap[repoChoice].Path)
	} else {
		fmt.Printf("[Err] Repository %s not found in the config\n", repoChoice)
		os.Exit(1)
	}

}