package configHelpers

import (
	"fmt"
	"os"
)

const configFileName = ".grconfig"

func ReadConfig() *os.File {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("[Err] Unable to read home directory\n", err)
		os.Exit(1)
	}

	fmt.Printf("[Info] Reading config file located in %s/%s\n", homeDir, configFileName)

	configFileObject, err := os.OpenFile(fmt.Sprintf("%s/%s", homeDir, configFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("[Err] Unable to read/create config file\n", err)
		os.Exit(1)
	}

	return configFileObject
}

func CloseConfig(configObject *os.File) {
	fmt.Println("[Info] Closing config file", configObject.Name())
	err := configObject.Close()
	if err != nil {
		fmt.Println("[Err] Unable to close file\n", err)
		os.Exit(1)
	}
}
