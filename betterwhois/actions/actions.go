package actions

import (
	"encoding/json"
	"fmt"
	"github.com/likexian/whois"
	"os"
	"strings"
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
	Name         string
	Organization string
	Country      string
	State        string
	City         string
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

func ParseData(rawWhoIs string) Domain {
	var (
		domainStruct                                                                      Domain
		line                                                                              string
		DomainName, DomainExpired, DomainCreated, DNSSEC                                  string
		RegistrarName, RegistrarEmail, RegistrarPhone, RegistrarURL, RegistrarWhois       string
		RegistrantName, RegistrantOrg, RegistrantCountry, RegistrantState, RegistrantCity string
		lineContent                                                                       []string
	)

	for _, line = range strings.Split(rawWhoIs, "\n") {
		line = strings.Trim(line, " ")
		lineContent = strings.Split(line, ":")
		// Skip matches that don't have the needed number of fields
		if len(lineContent) <= 1 {
			break
		}

		if lineContent[1] == "REDACTED FOR PRIVACY" {
			break
		}

		if len(lineContent) == 2 {
			lineContent[1] = strings.Trim(lineContent[1], " ")
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
		case "Registrant Name":
			RegistrantName = lineContent[1]
		case "Registrant Organization":
			RegistrantOrg = lineContent[1]
		case "Registrant Country":
			RegistrantCountry = lineContent[1]
		case "Registrant State/Province":
			RegistrantState = lineContent[1]
		case "Registrant City":
			RegistrantCity = lineContent[1]
		}

	}

	domainStruct = Domain{
		Name:    DomainName,
		Expiry:  DomainExpired,
		Created: DomainCreated,
		DNS: DNS{
			NameServers: []string{"Something1", "Something2"},
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
			Name:         RegistrantName,
			Organization: RegistrantOrg,
			Country:      RegistrantCountry,
			State:        RegistrantState,
			City:         RegistrantCity,
		},
	}

	return domainStruct
}

func (domain Domain) Print() {
	jsonOutput, _ := json.Marshal(domain)
	fmt.Println(string(jsonOutput))
}
