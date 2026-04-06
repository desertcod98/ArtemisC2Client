//go:build !debug

package dns

import (
	"fmt"
	"net"
)

var (
	DomainName = ".artemis.com."
)

// TXT query using system resolver
func DnsQuery(subdomain string) (string, error) {
	query := subdomain + DomainName
	records, err := net.LookupTXT(query)
	if err != nil {
		return "", err
	}
	if len(records) == 0 {
		return "", fmt.Errorf("no TXT record found for %s", query)
	}
	return records[0], nil
}
