package actions

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/likexian/whois"
	"os"
	"slices"
	"strings"
	"unicode"
)

var (
	Yellow = color.New(color.FgYellow).SprintFunc()
	Cyan   = color.New(color.FgCyan).SprintFunc()
	Green  = color.New(color.FgGreen).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
)

type Domain struct {
	Name       string `default:"Redacted for privacy"`
	Expiry     string `default:"Redacted for privacy"`
	Created    string `default:"Redacted for privacy"`
	Registrar  Registrar
	Registrant Registrant
	DNS        DNS
}

type DNS struct {
	NameServers []string
	DNSSEC      string `default:"Redacted for privacy"`
}

type Registrar struct {
	Name  string `default:"Redacted for privacy"`
	Email string `default:"Redacted for privacy"`
	Phone string `default:"Redacted for privacy"`
	URL   string `default:"Redacted for privacy"`
	Whois string `default:"Redacted for privacy"`
}

type Registrant struct {
	Organization string `default:"Redacted for privacy"`
	Country      string `default:"Redacted for privacy"`
	State        string `default:"Redacted for privacy"`
}

type Ip struct {
	Value        string
	NetInfo      NetInfo
	Organization string `default:"Redacted for privacy"`
	Country      string `default:"Redacted for privacy"`
	State        string `default:"Redacted for privacy"`
	City         string `default:"Redacted for privacy"`
	Address      string `default:"Redacted for privacy"`
	PostalCode   string `default:"Redacted for privacy"`
}

type NetInfo struct {
	Netrange string `default:"Redacted for privacy"`
	Netname  string `default:"Redacted for privacy"`
}

func GetType(userInput string) string {
	var (
		inputType                    string
		character                    int32
		dotCounter, characterCounter int
	)

	for _, character = range userInput {
		// Check for dot characters (.)(.)
		if character == 46 {
			dotCounter += 1
		}

		if !unicode.IsDigit(character) && character != 46 {
			characterCounter += 1
		}
	}

	switch dotCounter {
	case 0:
		fmt.Printf("[Err] The input provided %s is incorrect\n", userInput)
		os.Exit(1)
	case 1:
		inputType = "domain"
	case 2:
		inputType = "domain"
	default:
		if characterCounter == 0 {
			inputType = "ip"
		} else {
			fmt.Printf("[Err] The input provided %s is incorrect\n", userInput)
			os.Exit(1)
		}

	}

	return inputType
}

func GetWhois(userInput string) string {
	var (
		rawWhoIs string
		err      error
	)

	rawWhoIs, err = whois.Whois(userInput)
	if err != nil {
		fmt.Printf("[Err] Unable to get whois data for %s\n%s", userInput, err)
		os.Exit(1)
	}

	if strings.Contains(rawWhoIs, "No match for") {
		fmt.Printf("[Err] Unable to get whois data for %s\n", userInput)
		os.Exit(1)
	}
	return rawWhoIs
}

func ParseIpData(rawWhoIs string, inputIp string) Ip {
	var (
		ipStruct                                                                         Ip
		sliceCounter                                                                     int
		line, NetRange, NetName, Organization, Country, State, City, Address, PostalCode string
		lineContent                                                                      []string
	)

	for _, line = range strings.Split(rawWhoIs, "\n") {
		lineContent = strings.Split(line, ":")

		// Trim whitespace from results
		if len(lineContent) == 2 {
			lineContent[1] = strings.Trim(lineContent[1], " ")
		}

		// Skip matches that don't have the needed number of fields
		if len(lineContent) <= 1 {
			continue
		}

		// Skip redacted lines
		if lineContent[1] == "REDACTED FOR PRIVACY" {
			continue
		}

		for sliceCounter = range lineContent {
			lineContent[sliceCounter] = strings.TrimSpace(lineContent[sliceCounter])
		}

		switch lineContent[0] {
		case "NetRange":
			NetRange = strings.Join(lineContent[1:len(lineContent)], "")
		case "NetName":
			NetName = lineContent[1]
		case "OrgName":
			Organization = lineContent[1]
		case "Country":
			Country = lineContent[1]
		case "StateProv":
			State = lineContent[1]
		case "City":
			City = lineContent[1]
		case "Address":
			Address = lineContent[1]
		case "PostalCode":
			PostalCode = lineContent[1]
		}
	}

	ipStruct = Ip{
		Value: inputIp,
		NetInfo: NetInfo{
			Netrange: NetRange,
			Netname:  NetName,
		},
		Organization: Organization,
		Country:      Country,
		State:        State,
		City:         City,
		Address:      Address,
		PostalCode:   PostalCode,
	}

	return ipStruct
}

func ParseDomainData(rawWhoIs string) Domain {
	var (
		domainStruct                                                                Domain
		sliceCounter                                                                int
		line                                                                        string
		DomainName, DomainExpired, DomainCreated, DNSSEC                            string
		RegistrarName, RegistrarEmail, RegistrarPhone, RegistrarURL, RegistrarWhois string
		RegistrantOrg, RegistrantCountry, RegistrantState                           string
		lineContent, Nameservers                                                    []string
	)

	for _, line = range strings.Split(rawWhoIs, "\n") {
		lineContent = strings.Split(line, ":")

		// Trim whitespace from results
		for sliceCounter = range lineContent {
			lineContent[sliceCounter] = strings.TrimSpace(lineContent[sliceCounter])
		}

		// Skip matches that don't have the needed number of fields
		if len(lineContent) <= 1 {
			continue
		}

		// Skip redacted lines
		if lineContent[1] == "REDACTED FOR PRIVACY" {
			continue
		}

		switch lineContent[0] {
		case "Domain Name":
			DomainName = strings.ToLower(lineContent[1])
		case "Registry Expiry Date":
			DomainExpired = lineContent[1]
		case "Creation Date":
			DomainCreated = lineContent[1]
		case "DNSSEC":
			DNSSEC = lineContent[1]
		case "Name Server":
			Nameservers = append(Nameservers, strings.ToLower(lineContent[1]))
		case "Registrar":
			RegistrarName = lineContent[1]
		case "Registrar Abuse Contact Email":
			RegistrarEmail = lineContent[1]
		case "Registrar Abuse Contact Phone":
			RegistrarPhone = lineContent[1]
		case "Registrar URL":
			RegistrarURL = strings.Join(lineContent[1:len(lineContent)], "")
		case "Registrar WHOIS Server":
			RegistrarWhois = lineContent[1]
		case "Registrant Organization":
			RegistrantOrg = lineContent[1]
		case "Registrant Country":
			RegistrantCountry = lineContent[1]
		case "Registrant State/Province":
			RegistrantState = lineContent[1]
		}

	}

	// Remove duplicates from the NameServers slice
	slices.Sort(Nameservers)
	Nameservers = slices.Compact(Nameservers)

	domainStruct = Domain{
		Name:    DomainName,
		Expiry:  DomainExpired,
		Created: DomainCreated,
		DNS: DNS{
			NameServers: Nameservers,
			DNSSEC:      DNSSEC,
		},
		Registrar: Registrar{
			Name:  RegistrarName,
			Email: RegistrarEmail,
			Phone: RegistrarPhone,
			URL:   RegistrarURL,
			Whois: RegistrarWhois,
		},
		Registrant: Registrant{
			Organization: RegistrantOrg,
			Country:      RegistrantCountry,
			State:        RegistrantState,
		},
	}

	return domainStruct
}

func (ip Ip) Print() {
	fmt.Printf("%s: %s\n\n", colorPrint("Checking IP", "title"), colorPrint(ip.Value, "header"))
	fmt.Printf("%s: %s\n", colorPrint("Organization", "title"), colorPrint(ip.Organization, "content"))
	fmt.Printf("%s: %s\n", colorPrint("Country", "title"), colorPrint(ip.Country, "content"))
	fmt.Printf("%s: %s\n", colorPrint("State", "title"), colorPrint(ip.State, "content"))
	fmt.Printf("%s: %s\n", colorPrint("City", "title"), colorPrint(ip.City, "content"))
	fmt.Printf("%s: %s\n", colorPrint("Address", "title"), colorPrint(ip.Address, "content"))
	fmt.Printf("%s: %s\n", colorPrint("PostalCode", "title"), colorPrint(ip.PostalCode, "content"))
	fmt.Printf("%s:\n", colorPrint("Netinfo", "title"))
	fmt.Printf("  %s: %s\n", colorPrint("Netrange", "title"), colorPrint(ip.NetInfo.Netrange, "content"))
	fmt.Printf("  %s: %s\n", colorPrint("Netname", "title"), colorPrint(ip.NetInfo.Netname, "content"))

}

func (domain Domain) Print() {
	fmt.Printf("%s: %s\n\n", colorPrint("Checking Domain", "title"), colorPrint(domain.Name, "header"))
	fmt.Printf("%s: %s\n", colorPrint("Creation date", "title"), colorPrint(domain.Created, "content"))
	fmt.Printf("%s: %s\n", colorPrint("Expiration date", "title"), colorPrint(domain.Expiry, "content"))
	fmt.Printf("%s:\n", colorPrint("DNS info", "title"))
	fmt.Printf("  %s:\n", colorPrint("NameServers", "title"))
	for _, arrayElement := range domain.DNS.NameServers {
		fmt.Printf("    - %s\n", colorPrint(arrayElement, "content"))
	}
	fmt.Printf("  %s: %s\n", colorPrint("DNSSEC", "title"), colorPrint(domain.DNS.DNSSEC, "content"))
	fmt.Printf("%s:\n", colorPrint("Registrant info", "title"))
	fmt.Printf("  %s: %s\n", colorPrint("Organization", "title"), colorPrint(domain.Registrant.Organization, "content"))
	fmt.Printf("  %s: %s\n", colorPrint("Country", "title"), colorPrint(domain.Registrant.Country, "content"))
	fmt.Printf("  %s: %s\n", colorPrint("State", "title"), colorPrint(domain.Registrant.State, "content"))
	fmt.Printf("%s:\n", colorPrint("Registrar info", "title"))
	fmt.Printf("  %s: %s\n", colorPrint("Name", "title"), colorPrint(domain.Registrar.Name, "content"))
	fmt.Printf("  %s: %s\n", colorPrint("Email", "title"), colorPrint(domain.Registrar.Email, "content"))
	fmt.Printf("  %s: %s\n", colorPrint("Phone", "title"), colorPrint(domain.Registrar.Phone, "content"))
	fmt.Printf("  %s: %s\n", colorPrint("Url", "title"), colorPrint(domain.Registrar.URL, "content"))
	fmt.Printf("  %s: %s\n", colorPrint("WhoIs", "title"), colorPrint(domain.Registrar.Whois, "content"))

}

func colorPrint(inputText string, textType string) string {
	if inputText == "Redacted for privacy" {
		return Red(inputText)
	}

	switch textType {
	case "header":
		return Yellow(inputText)
	case "title":
		return Cyan(inputText)
	case "content":
		return Green(inputText)
	default:
		fmt.Printf("Unknown textType: %s\n", textType)
		os.Exit(1)
	}

	return "empty"
}
