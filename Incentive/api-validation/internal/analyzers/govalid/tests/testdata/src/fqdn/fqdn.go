//go:generate ./fqdn.go
package fqdn

// Server is a struct for testing fqdn validation
type Server struct {
	// +govalid:fqdn
	Hostname string `json:"hostname"`

	// +govalid:fqdn
	Domain string `json:"domain"`

	// +govalid:fqdn
	MailServer string `json:"mail_server"`
}
