package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"os"
)

const logFileName = "fsync.log"

type InputArgs struct {
	Action     string
	ConfigFile os.File
	PublicKey  ssh.Signer
	Hosts      ssh.HostKeyCallback
	LogFile    os.File // Not in use yet
}

type HostConfig struct {
	HostsMap map[string]struct {
		Hostname  string `json:"hostname"`
		Port      int    `json:"port"`
		User      string `json:"user"`
		LocalDir  string `json:"local_dir"`
		RemoteDir string `json:"remote_dir"`
	}
	SSHKey ssh.Signer
	Hosts  ssh.HostKeyCallback
	Logger string // Not in use yet
}

func ArgInit() InputArgs {
	argParser := argparse.NewParser("fsync", "File synchronisation service for code editors")
	selectedAction := argParser.StringPositional(&argparse.Options{Help: "Action which should be performed", Default: "run"})
	configFile := argParser.File("f", "file", os.O_RDWR, 0644, &argparse.Options{Required: true, Help: "Location of config file"})
	sshKey := argParser.File("k", "key", os.O_RDONLY, 0644, &argparse.Options{Required: true, Help: "Location of the private key"})
	hostsFile := argParser.File("j", "hosts", os.O_RDONLY, 0644, &argparse.Options{Required: true, Help: "Location of the hosts file"})
	logFile := argParser.File("l", "log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644, &argparse.Options{Required: false, Help: "Location of file for logging", Default: logFileName})

	err := argParser.Parse(os.Args)
	if err != nil {
		fmt.Print(argParser.Usage(err))
	}

	keyData, err := os.ReadFile(sshKey.Name())
	if err != nil {
		fmt.Println("Error reading private key:", err)
		os.Exit(1)
	}

	privateKey, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		os.Exit(1)
	}

	hostsData, err := knownhosts.New(hostsFile.Name())
	if err != nil {
		fmt.Println("Error parsing hosts file:", err)
		os.Exit(1)
	}

	return InputArgs{*selectedAction, *configFile, privateKey, hostsData, *logFile}
}

func BuildHostConfig(i InputArgs) HostConfig {
	var hosts HostConfig

	data, err := os.ReadFile(i.ConfigFile.Name())
	if err != nil {
		fmt.Println("Encountered error while reading json:", err)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &hosts.HostsMap)
	if err != nil {
		fmt.Println("Encountered error while unmarshalling file:", err)
		os.Exit(1)
	}

	for key, value := range hosts.HostsMap {
		if value.Port == 0 {
			value.Port = 22

			hosts.HostsMap[key] = value
		}
	}

	hosts.SSHKey = i.PublicKey
	hosts.Hosts = i.Hosts
	hosts.Logger = i.LogFile.Name()

	return hosts
}

func (hosts HostConfig) VerifyHosts() {
	for hostPetName, hostData := range hosts.HostsMap {
		fmt.Printf("[%s] Starting verification\n", hostPetName)

		sshConfig := &ssh.ClientConfig{
			User: hostData.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(hosts.SSHKey),
			},
			HostKeyCallback: hosts.Hosts,
		}

		conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostData.Hostname, hostData.Port), sshConfig)
		if err != nil {
			fmt.Println("Encountered error trying to connect over ssh:", err)
			os.Exit(1)
		}

		client, err := sftp.NewClient(conn)
		if err != nil {
			fmt.Println("Encountered error while trying to create client:", err)
			os.Exit(1)
		}

		defer client.Close()
		pwd, err := client.Getwd()
		_ = pwd
		if err != nil {
			fmt.Println("Encountered error while reading directory:", err)
			os.Exit(1)
		} else {
			fmt.Printf("[%s] Verification succesfull\n", hostPetName)
		}
	}
}
