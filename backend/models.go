package models

type Subdomain struct {
	Name       string
	IPs        []string
	IPOwner    string
	CertIssuer string
	CertExpiry string
}
