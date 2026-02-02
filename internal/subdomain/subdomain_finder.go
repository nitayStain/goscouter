// Subdomain discovery via Certificate Transparency and DNS resolution.

package subdomain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	defaultUserAgent   = "goscouter-subdomain-finder/1.0"
	defaultMaxBodySize = 25 << 20
)

var ErrInvalidDomain = errors.New("invalid domain")

type Finder struct {
	httpClient     *http.Client
	resolver       *net.Resolver
	userAgent      string
	maxBodySize    int64
	lookupIPOwners bool
	debug          bool
}

type FinderOption func(*Finder)

func NewFinder(opts ...FinderOption) *Finder {
	f := &Finder{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		resolver:       net.DefaultResolver,
		userAgent:      defaultUserAgent,
		maxBodySize:    defaultMaxBodySize,
		lookupIPOwners: true,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(f)
		}
	}
	return f
}

func WithHTTPClient(client *http.Client) FinderOption {
	return func(f *Finder) {
		if client != nil {
			f.httpClient = client
		}
	}
}

func WithResolver(resolver *net.Resolver) FinderOption {
	return func(f *Finder) {
		if resolver != nil {
			f.resolver = resolver
		}
	}
}

func WithUserAgent(userAgent string) FinderOption {
	return func(f *Finder) {
		userAgent = strings.TrimSpace(userAgent)
		if userAgent != "" {
			f.userAgent = userAgent
		}
	}
}

func WithMaxBodySize(maxBytes int64) FinderOption {
	return func(f *Finder) {
		if maxBytes > 0 {
			f.maxBodySize = maxBytes
		}
	}
}

func WithIPOwnerLookup(enabled bool) FinderOption {
	return func(f *Finder) {
		f.lookupIPOwners = enabled
	}
}

func WithDebug(enabled bool) FinderOption {
	return func(f *Finder) {
		f.debug = enabled
	}
}

func SubdomainFinder(domain string) (map[string]Subdomain, bool, error) {
	return NewFinder().Find(context.Background(), domain)
}

type crtShEntry struct {
	NameValue  string `json:"name_value"`
	IssuerName string `json:"issuer_name,omitempty"`
	NotAfter   string `json:"not_after,omitempty"`
}

type certSpotterEntry struct {
	DNSNames []string `json:"dns_names"`
	Issuer   string   `json:"issuer,omitempty"`
	NotAfter string   `json:"not_after,omitempty"`
}

func (f *Finder) Find(ctx context.Context, domain string) (map[string]Subdomain, bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	normalizedDomain, err := normalizeDomain(domain)
	if err != nil {
		return nil, false, err
	}

	results := make(map[string]Subdomain)
	var sourceErrs []error

	if err := f.collectFromCrtSh(ctx, normalizedDomain, results); err != nil {
		sourceErrs = append(sourceErrs, err)
	}
	if err := f.collectFromCertSpotter(ctx, normalizedDomain, results); err != nil {
		sourceErrs = append(sourceErrs, err)
	}

	if len(results) == 0 && len(sourceErrs) > 0 {
		return nil, false, errors.Join(sourceErrs...)
	}

	hasWildcard, _ := f.hasWildcardDNS(ctx, normalizedDomain)
	f.enrichResults(ctx, results)

	return results, hasWildcard, nil
}

func (f *Finder) collectFromCrtSh(ctx context.Context, domain string, results map[string]Subdomain) error {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	var entries []crtShEntry
	if err := f.fetchJSON(ctx, url, &entries); err != nil {
		return fmt.Errorf("crt.sh request failed: %w", err)
	}
	for _, entry := range entries {
		for _, name := range strings.Split(entry.NameValue, "\n") {
			name = normalizeName(name)
			if !isSubdomainOf(name, domain) {
				continue
			}
			addSubdomain(results, name, entry.IssuerName, entry.NotAfter)
		}
	}
	return nil
}

func (f *Finder) collectFromCertSpotter(ctx context.Context, domain string, results map[string]Subdomain) error {
	url := fmt.Sprintf(
		"https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names",
		domain,
	)
	var entries []certSpotterEntry
	if err := f.fetchJSON(ctx, url, &entries); err != nil {
		return fmt.Errorf("certspotter request failed: %w", err)
	}
	for _, entry := range entries {
		for _, name := range entry.DNSNames {
			name = normalizeName(name)
			if !isSubdomainOf(name, domain) {
				continue
			}
			addSubdomain(results, name, entry.Issuer, entry.NotAfter)
		}
	}
	return nil
}

func (f *Finder) enrichResults(ctx context.Context, results map[string]Subdomain) {
	ownerCache := make(map[string]string)

	for name := range results {
		ips, err := f.resolveIPs(ctx, name)
		if err != nil || len(ips) == 0 {
			delete(results, name)
			continue
		}

		data := results[name]
		data.IPs = ips

		if f.lookupIPOwners {
			owners := make([]string, 0, len(ips))
			for _, ip := range ips {
				owner, ok := ownerCache[ip]
				if !ok {
					owner, _ = f.getIPOwner(ctx, ip)
					ownerCache[ip] = owner
				}
				if owner != "" {
					owners = append(owners, owner)
				}
			}
			if len(owners) > 0 {
				data.IPOwner = strings.Join(uniqueStrings(owners), ", ")
			}
		}

		results[name] = data
	}
}

func (f *Finder) hasWildcardDNS(ctx context.Context, domain string) (bool, error) {
	sub, err := randomSubdomain()
	if err != nil {
		return false, err
	}
	randomSub := fmt.Sprintf("%s.%s", sub, domain)

	if f.debug {
		log.Printf("[DEBUG] Checking wildcard DNS: %s", randomSub)
	}

	_, err = f.resolver.LookupHost(ctx, randomSub)
	hasWildcard := err == nil

	if f.debug {
		log.Printf("[DEBUG] Wildcard DNS result for %s: %v", domain, hasWildcard)
	}

	return hasWildcard, nil
}

func (f *Finder) resolveIPs(ctx context.Context, domain string) ([]string, error) {
	if f.debug {
		log.Printf("[DEBUG] DNS lookup: %s", domain)
	}

	ips, err := f.resolver.LookupHost(ctx, domain)
	if err != nil {
		if f.debug {
			log.Printf("[DEBUG] DNS lookup failed for %s: %v", domain, err)
		}
		return nil, err
	}

	uniqueIPs := uniqueStrings(ips)
	if f.debug {
		log.Printf("[DEBUG] DNS resolved %s -> %v", domain, uniqueIPs)
	}

	return uniqueIPs, nil
}

func (f *Finder) fetchJSON(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", f.userAgent)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body := io.LimitReader(resp.Body, f.maxBodySize)
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(out); err != nil {
		return err
	}
	return nil
}

func (f *Finder) fetchText(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", f.userAgent)

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body := io.LimitReader(resp.Body, 2048)
	data, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func (f *Finder) getIPOwner(ctx context.Context, ip string) (string, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s/org", ip)
	owner, err := f.fetchText(ctx, url)
	if err != nil {
		return "", err
	}
	return owner, nil
}

func normalizeDomain(domain string) (string, error) {
	domain = strings.TrimSpace(strings.ToLower(domain))
	domain = strings.TrimSuffix(domain, ".")
	if domain == "" || !isValidDomain(domain) {
		return "", fmt.Errorf("%w: %q", ErrInvalidDomain, domain)
	}
	return domain, nil
}

func normalizeName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.TrimSuffix(name, ".")
	name = strings.TrimPrefix(name, "*.")
	if strings.Contains(name, "*") {
		return ""
	}
	return name
}

func isSubdomainOf(name, domain string) bool {
	if name == "" {
		return false
	}
	if name == domain {
		return true
	}
	if !strings.HasSuffix(name, "."+domain) {
		return false
	}
	return isValidDomain(name)
}

func addSubdomain(results map[string]Subdomain, name, issuer, expiry string) {
	name = normalizeName(name)
	if name == "" {
		return
	}

	entry, exists := results[name]
	if !exists {
		results[name] = Subdomain{
			Name:       name,
			CertIssuer: issuer,
			CertExpiry: expiry,
		}
		return
	}

	if entry.CertIssuer == "" && issuer != "" {
		entry.CertIssuer = issuer
	}
	if entry.CertExpiry == "" && expiry != "" {
		entry.CertExpiry = expiry
	}
	results[name] = entry
}

func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func isValidDomain(name string) bool {
	if len(name) == 0 || len(name) > 253 {
		return false
	}

	labels := strings.Split(name, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		if label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for i := 0; i < len(label); i++ {
			ch := label[i]
			if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
				continue
			}
			return false
		}
	}
	return true
}

func randomSubdomain() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
