// POC of a subdomain finder via Certificate Transparency

package models

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type CrtShEntry struct {
	NameValue  string `json:"name_value"`
	IssuerName string `json:"issuer_name,omitempty"`
	NotAfter   string `json:"not_after,omitempty"`
}

type CertSpotterEntry struct {
	DNSNames []string `json:"dns_names"`
	Issuer   string   `json:"issuer,omitempty"`
	NotAfter string   `json:"not_after,omitempty"`
}

func SubdomainFinder(domain string) (map[string]Subdomain, bool, error) {
	unique := make(map[string]Subdomain)

	// quarry subdomains from crt.sh
	crtURL := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	if body, err := fetch(crtURL); err == nil {
		var entries []CrtShEntry
		if json.Unmarshal(body, &entries) == nil {
			for _, entry := range entries {
				for _, name := range strings.Split(entry.NameValue, "\n") {
					name = strings.TrimSpace(name)
					if name != "" {
						unique[name] = Subdomain{
							Name:       name,
							CertIssuer: entry.IssuerName,
							CertExpiry: entry.NotAfter,
						}
					}
				}
			}
		}
	}

	// quarry subdomains from certspotter.com
	csURL := fmt.Sprintf(
		"https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names",
		domain,
	)
	if body, err := fetch(csURL); err == nil {
		var entries []CertSpotterEntry
		if json.Unmarshal(body, &entries) == nil {
			for _, entry := range entries {
				for _, name := range entry.DNSNames {
					name = strings.TrimSpace(name)
					if name != "" {
						if _, exists := unique[name]; !exists {
							unique[name] = Subdomain{
								Name:       name,
								CertIssuer: entry.Issuer,
								CertExpiry: entry.NotAfter,
							}
						}
					}
				}
			}
		}
	}

	hasWildcard, err := hasWildcardDNS(domain)
	if err != nil {
		return nil, false, err
	}

	for name, data := range unique {
		ips, err := resolveIPs(name)
		if err != nil || len(ips) == 0 {
			delete(unique, name)
			continue
		}
		data.IPs = ips

		owners := []string{}
		for _, ip := range ips {
			if owner, err := getHost(ip); err == nil && owner != "" {
				owners = append(owners, owner)
			}
		}
		if len(owners) > 0 {
			data.IPOwner = strings.Join(owners, ", ")
		}

		unique[name] = data
	}

	return unique, hasWildcard, nil
}

// *---TOOLS---* //

// Chuck if the domain has a Wildcard Record
func hasWildcardDNS(domain string) (bool, error) {
	sub, err := randomSubdomain()
	if err != nil {
		return false, err
	}
	randomSub := fmt.Sprintf("%s.%s", sub, domain)
	_, err = net.LookupHost(randomSub)
	return err == nil, nil
}

// Resolve IP addresses for a subdomain
func resolveIPs(domain string) ([]string, error) {
	return net.LookupHost(domain)
}

// Fetch helper
func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// Get IP range owner from ipinfo.io
func getHost(ip string) (string, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s/org", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}

// generate a random subdomain for the wildcard check
func randomSubdomain() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
