package dns

import (
	"fmt"
	"github.com/miekg/dns"
)

var client = dns.Client{}

const dnsServer = "127.0.0.1:53"

func DnsQuery(subdomain string) (string, error) {
	m := dns.Msg{}
	query := subdomain + ".artemis.com."
	m.SetQuestion(query, dns.TypeTXT)

	r, _, err := client.Exchange(&m, dnsServer)
	if err != nil {
		return "", err
	}

	for _, ans := range r.Answer {
		t, ok := ans.(*dns.TXT)
		if ok && len(t.Txt) > 0 {
			return t.Txt[0], nil
		}
	}
	return "", fmt.Errorf("no TXT record found for %s", query)
}
