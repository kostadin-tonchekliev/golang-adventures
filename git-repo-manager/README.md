Golang Repository Manager
---
version 0.1<br>

## Installation:
### 1. Pre-built
The pre-built installation is straight forward. Copy the necessary executable from the `executables` folder to your desired location and run it. Ideally this needs to be in `/usr/local/bin/` and renamed to without the OS info (`gmanager`)<br>
Once the executable is installed just run `gmanager setup` to configure the project environment. (assuming your executable is called `gmanager`)

### 2. Manual building
In order to manually build the executable for your environment you need to have [Golang](https://go.dev/) installed. Once this is done clone the repository and build the executable using `go build gmanager.go`. Once the executable is build move it to your desired location, ideally this needs to be in `/usr/local/bin`<br>
Similar to the pre-built steps run `gmanager setup` to configure the project environment. (assuming your executable is called `gmanager`).

## Usage:
In order to do anything with the script you will need to build a config with the repositories, for this run the following command which will prompt you to add the necessery information:
```shell
gmanager config add
```
Once the config is build you can run any of the actions (see supported actions below for more info)

## Supported Actions:
* `setup` - Configure project environment. This needs to be run only on the first run.
* `status` - Show the status of the repositories. Prints repository petname, current local branch and types/numbers of staged files.
* `version` - Display build information
* `config` - Actions for modifying the config
  * `ls` - List existing config
  * `add`- Add new options to the config
  * `remove` - Remove options from the config
* `cd` - Move to repositories. **This isn't to be run on its own!**
