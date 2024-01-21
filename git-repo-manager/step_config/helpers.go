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

func readInput(message string) string {
	var inputString string

	for {
		fmt.Printf("[Input] %s: ", message)
		_, err := fmt.Scanf("%s", &inputString)
		if err != nil {
			fmt.Println("[Err] Empty value provided!")
		} else {
			break
		}
	}

	return inputString
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

	// Check if the file is empty which means that it didn't exist beforehand
	configFileStat, _ := configFileObject.Stat()
	if configFileStat.Size() == 0 {
		configFileObject.Write([]byte(`{}`))
		configFileObject.Close()
		configFileObject, err = os.OpenFile(fmt.Sprintf("%s/%s", homeDir, configFileName), os.O_RDWR, 0644)
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
	fmt.Println("[Info] Closing config file", config.ConfigFile.Name())
	err := config.ConfigFile.Close()
	if err != nil {
		fmt.Println("[Err] Unable to close file\n", err)
		os.Exit(1)
	}
}

func (config Config) ListConfig() {
	fmt.Println("[Info] Reading content of config", config.ConfigFile.Name())
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
	petName = readInput("Select petname for the repository")
	uri = readInput("Select uri of the repository")

	for {
		path = readInput("Select local path of the repository")
		files, err := os.ReadDir(path)
		if err != nil || len(files) == 0 {
			fmt.Println("[Err] Please select valid directory")
		} else {
			break
		}
	}

	fmt.Printf("[Debug] New information: %s - %s - %s\n", petName, uri, path)
	config.RepoMap[petName] = RepoObject{
		Url:  uri,
		Path: path,
	}

	jsonContent, err := json.MarshalIndent(config.RepoMap, "", "  ")
	if err != nil {
		fmt.Println("[Err] Unable to convert map to json\n", err)
	}

	fmt.Println(jsonContent)

	config.ConfigFile.Truncate(0)
	_, err = config.ConfigFile.Write(jsonContent)
	if err != nil {
		fmt.Println("[Err] Unable to write new value to config\n", err)
	}
}
