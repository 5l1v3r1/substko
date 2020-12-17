package substko

import (
	"bytes"
	"errors"
	"strings"

	"github.com/drsigned/gos"
	"github.com/haccer/available"
)

// Fingerprint is a
type Fingerprint struct {
	Service     string   `json:"service"`
	Status      string   `json:"status"`
	Cname       []string `json:"cname"`
	Fingerprint []string `json:"fingerprint"`
}

// Options is a
type Options struct {
	HTTPS        bool
	Timeout      int
	Fingerprints []Fingerprint
}

// CheckNSSTKO is a
func CheckNSSTKO(target string) (status string, STKOType string, at string, err error) {
	// 1. Resolve all nameservers for the domain
	nonAuthoritativeNS, err := nslookup(target, "8.8.8.8")
	if err != nil {
		return "Not Vulnerable", "NS", "", err
	}

	if len(nonAuthoritativeNS) == 0 {
		return "Not Vulnerable", "NS", "", errors.New("couldn't find any non authoritative nameservers for " + target)
	}

	authoritativeNS, err := nslookup(target, nonAuthoritativeNS[0])
	if err != nil {
		return "Not Vulnerable", "NS", "", err

	}

	for _, nameserver := range authoritativeNS {
		// 2. FOR CUSTOM NAMESERVER DOMAINS: Check if nameserver base domain is available for registration.

		u, err := gos.ParseURL(unFQDN(nameserver))
		if err != nil {
			return "Not Vulnerable", "NS", "", err
		}

		if hostLookup(unFQDN(nameserver)) && available.Domain(u.ETLDPlus1) {
			return "Vulnerable", "NS", nameserver, nil
		}

		// 3. FOR MANAGED DNS: resolve domain against this nameserver and look for responses like SERVFAIL or REFUSED.
		ns, err := nsReturnsWeirdStatus(target, nameserver)
		if err != nil {
			return "Not Vulnerable", "NS", "", err

		}

		if ns {
			return "Vulnerable", "NS", nameserver, nil
		}
	}

	return "Not Vulnerable", "NS", "", nil
}

// CheckCNAMESTKO is a
func CheckCNAMESTKO(subdomain string, o *Options) (status string, STKOType string, at string, err error) {
	cname, err := cnameLookup(subdomain)
	if err != nil {
		return "Not Vulnerable", "CNAME", "", err
	} else if cname == "" {
		return "Not Vulnerable", "CNAME", "", nil
	}

	u, err := gos.ParseURL(unFQDN(cname))
	if err != nil {
		return "Not Vulnerable", "CNAME", "", err
	}

	if hostLookup(unFQDN(cname)) && available.Domain(u.ETLDPlus1) {

		return "Vulnerable", "CNAME", unFQDN(cname), nil
	}

	serviceIndex := -1

CNAME:
	for i := range o.Fingerprints {
		for j := range o.Fingerprints[i].Cname {
			if strings.Contains(unFQDN(cname), o.Fingerprints[i].Cname[j]) {
				serviceIndex = i
				break CNAME
			}
		}
	}

	if serviceIndex >= 0 {
		body := getBody(subdomain, o.HTTPS, o.Timeout)

		for i := range o.Fingerprints[serviceIndex].Fingerprint {
			if bytes.Contains(body, []byte(o.Fingerprints[serviceIndex].Fingerprint[i])) {
				return o.Fingerprints[serviceIndex].Status, "CNAME", unFQDN(cname), nil
			}
		}
	}

	return "Not Vulnerable", "CNAME", "", nil
}

// CheckSTKO is a
func CheckSTKO(subdomain string, o *Options) (status string, STKOType string, at string, err error) {
	status, STKOType, at, err = CheckNSSTKO(subdomain)
	if status == "Not Vulnerable" || status == "" {
		status, STKOType, at, err = CheckCNAMESTKO(subdomain, o)
	}

	return status, STKOType, at, err
}
