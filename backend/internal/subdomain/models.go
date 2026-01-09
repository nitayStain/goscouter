package subdomain

// Subdomain captures data discovered for a subdomain name.
type Subdomain struct {
	Name       string   `json:"name"`
	IPs        []string `json:"ips,omitempty"`
	IPOwner    string   `json:"ip_owner,omitempty"`
	CertIssuer string   `json:"cert_issuer,omitempty"`
	CertExpiry string   `json:"cert_expiry,omitempty"`
}
