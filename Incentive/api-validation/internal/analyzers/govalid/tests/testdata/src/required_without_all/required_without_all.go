//go:generate ./required_without_all.go
package required_without_all

// Auth is a struct for testing required_without_all validation
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`

	// +govalid:required_without_all=Username Password
	APIKey string `json:"api_key"`

	SSHKey     string `json:"ssh_key"`
	SSHKeyPath string `json:"ssh_key_path"`

	// +govalid:required_without_all=SSHKey SSHKeyPath
	SSHPassword string `json:"ssh_password"`
}
