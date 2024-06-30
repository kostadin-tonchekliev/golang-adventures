package main

import (
	"betterwhois/actions"
	"fmt"
	"os"
)

func main() {
	var (
		sysArguments []string
		rawWhois     string
		domain       actions.Domain
	)

	sysArguments = os.Args
	switch {
	case len(sysArguments) == 1:
		fmt.Println("[Err] No arguments provided. Please select domain or IP")
		os.Exit(1)
	case len(sysArguments) == 2:
		rawWhois = actions.GetWhois(sysArguments[1])
		domain = actions.ParseData(rawWhois)
		domain.Print()
		os.Exit(0)
	case len(sysArguments) >= 3:
		fmt.Println("[Err] Too many arguments provided")
		os.Exit(1)
	}
}
