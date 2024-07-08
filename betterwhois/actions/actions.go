package actions

import (
	"encoding/json"
	"fmt"
	"github.com/likexian/whois"
	"os"
	"slices"
	"strings"
	"unicode"
)

type Domain struct {
	Name       string
	Expiry     string
	Created    string
	Registrar  Registrar
	Registrant Registrant
	DNS        DNS
}

type DNS struct {
	NameServers []string
	DNSSEC      string
}

type Registrar struct {
	Name  string
	Email string
	Phone string
	URL   string
	Whois string
}

type Registrant struct {
	Organization string
	Country      string
	State        string
}

type Ip struct {
	NetInfo      NetInfo
	Organization string
	Country      string
	State        string
	City         string
	Address      string
	PostalCode   string
}

type NetInfo struct {
	Netrange string
	Netname  string
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

	return rawWhoIs
}

func ParseIpData(rawWhoIs string) Ip {
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
			DomainName = lineContent[1]
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
	jsonOutput, _ := json.Marshal(ip)
	fmt.Println(string(jsonOutput))
}

func (domain Domain) Print() {
	jsonOutput, _ := json.Marshal(domain)
	fmt.Println(string(jsonOutput))
}
