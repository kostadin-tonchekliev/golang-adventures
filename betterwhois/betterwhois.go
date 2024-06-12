package main

import (
    "betterwhois/actions"
    "github.com/likexian/whois-parser"
)

func main() {
    var parsedData whoisparser.WhoisInfo
    parsedData = actions.ParseData("google.com")
    actions.DisplayData(parsedData)
}
