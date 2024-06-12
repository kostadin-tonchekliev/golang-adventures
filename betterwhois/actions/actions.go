package actions

import (
    "fmt"
    "os"
    "github.com/likexian/whois"
    "github.com/likexian/whois-parser"

)

func ParseData(domain string) whoisparser.WhoisInfo {
    var (
        rawWhoIs string
        parsedWhoIs whoisparser.WhoisInfo
        err error
    )

    rawWhoIs, err = whois.Whois(domain)
    if err != nil {
        fmt.Printf("[Err] Unable to get whois data for %s\n%s", domain, err)
        os.Exit(1)
    }

    parsedWhoIs, err = whoisparser.Parse(rawWhoIs)
    if err != nil {
        fmt.Printf("[Err] Unable to parse data for %s\n%s", domain, err)
        os.Exit(1)
    }

    return parsedWhoIs
}

func DisplayData(inputData whoisparser.WhoisInfo) {
    fmt.Println(inputData.Domain.Name)
}
