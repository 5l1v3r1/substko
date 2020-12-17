package substko

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
)

func unFQDN(domain string) string {
	return strings.TrimSuffix(domain, ".")
}

func parseNSRecords(records []dns.RR) (nameservers []string) {
	for _, ans := range records {
		if ans.Header().Rrtype == dns.TypeNS {
			record := ans.(*dns.NS)
			nameservers = append(nameservers, record.Ns)
		} else if ans.Header().Rrtype == dns.TypeSOA {
			record := ans.(*dns.SOA)
			nameservers = append(nameservers, record.Ns)
		}
	}

	return nameservers
}

func nslookup(domain string, nameserver string) (nameservers []string, err error) {
	client := dns.Client{}
	message := dns.Msg{}

	message.SetQuestion(dns.Fqdn(domain), dns.TypeNS)

	record, _, err := client.Exchange(&message, unFQDN(nameserver)+":53")
	if err != nil {
		return nil, err
	}

	if record.Rcode == dns.RcodeSuccess {
		if len(record.Answer) > 0 {
			nameservers = parseNSRecords(record.Answer)
		} else {
			// if no NS records are found in the answer section, fallback to using the authority section
			nameservers = parseNSRecords(record.Ns)
		}
	} else {
		return nil, errors.New("Failed!! to get NS servers")
	}

	return nameservers, nil
}

func nsReturnsWeirdStatus(subdomain string, nameserver string) (bool, error) {
	client := dns.Client{}
	message := dns.Msg{}

	message.SetQuestion(dns.Fqdn(subdomain), dns.TypeA)

	record, _, err := client.Exchange(&message, nameserver+":53")
	if err != nil {
		return false, err
	}

	if record.Rcode == dns.RcodeServerFailure || record.Rcode == dns.RcodeRefused {
		return true, nil
	}

	return false, nil
}

func cnameLookup(domain string) (cname string, err error) {
	client := dns.Client{}
	message := dns.Msg{}

	message.SetQuestion(dns.Fqdn(domain), dns.TypeCNAME)

	record, _, err := client.Exchange(&message, "8.8.8.8:53")
	if err != nil {
		return "", err
	} else if len(record.Answer) == 0 {
		return "", nil
	}

	recordData := record.Answer[len(record.Answer)-1].(*dns.CNAME)
	cname = recordData.Target

	for ok := true; ok; ok = len(record.Answer) > 0 {
		recordData = record.Answer[len(record.Answer)-1].(*dns.CNAME)
		cname = recordData.Target

		message.SetQuestion(dns.Fqdn(cname), dns.TypeCNAME)

		record, _, err = client.Exchange(&message, "8.8.8.8:53")
		if err != nil {
			break
		}
	}

	return cname, nil
}

func doesResolve(domain string) (bool, error) {
	domain = dns.Fqdn(domain)

	client := dns.Client{}
	message := dns.Msg{}

	message.SetQuestion(domain, dns.TypeA)

	record, _, err := client.Exchange(&message, "8.8.8.8:53")
	if err != nil {
		return false, err
	}

	if record.Rcode == dns.RcodeNameError {
		return false, nil
	}

	return true, nil
}

func hostLookup(domain string) bool {
	if _, err := net.LookupHost(domain); err != nil {
		if strings.Contains(fmt.Sprintln(err), "no such host") {
			return true
		}
	}

	return false
}
