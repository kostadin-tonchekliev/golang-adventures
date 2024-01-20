package configHelpers

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".grconfig.json"

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

	fmt.Printf("[Info] Reading config file located in %s/%s\n", homeDir, configFileName)

	configFileObject, err := os.OpenFile(fmt.Sprintf("%s/%s", homeDir, configFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("[Err] Unable to read/create config file\n", err)
		os.Exit(1)
	}

	configFileStat, _ := configFileObject.Stat()

	buffer := make([]byte, configFileStat.Size())
	configFileObject.Read(buffer) // Can probably add some error reporting here

	err = json.Unmarshal(buffer, &configFile.RepoMap)
	if err != nil {
		fmt.Println("[Err] Unable to unmarshal json\n", err)
		os.Exit(1)
	}

	configFile.ConfigFile = configFileObject

	return configFile
}

func (config Config) CloseConfig() {
	fmt.Println("[Info] Closing config file", config.ConfigFile.Name())
	err := config.ConfigFile.Close()
	if err != nil {
		fmt.Println("[Err] Unable to close file\n", err)
		os.Exit(1)
	}
}

func (config Config) ReadContent() {
	fmt.Println("[Info] Reading content of config", config.ConfigFile.Name())
	for petName, repoContent := range config.RepoMap {
		fmt.Println(petName)
		fmt.Println("\tURL:", repoContent.Url)
		fmt.Println("\tPath:", repoContent.Path)
		fmt.Println()
	}
}
