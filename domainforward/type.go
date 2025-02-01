package domainforward

import "github.com/mrehanabbasi/go-resellerclub/core"

type StdResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DetailsDomainForward struct {
	URLMasking          core.JSONBool `json:"urlmasking"`
	PathForwarding      core.JSONBool `json:"pathforwarding"`
	SubdomainForwarding core.JSONBool `json:"subdomainforwarding"`
	IPAddress           string        `json:"ipaddress"`
	DomainName          string        `json:"domainname"`
}

type DNSRecord struct {
	TimeToLive core.JSONInt `json:"timetolive"`
	Type       string       `json:"type"`
	Host       string       `json:"host"`
	Value      string       `json:"value"`
}
