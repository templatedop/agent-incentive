//go:generate ./excluded_without_all.go
package excluded_without_all

// System is a struct for testing excluded_without_all validation
type System struct {
	AdminUser     string `json:"admin_user"`
	AdminPassword string `json:"admin_password"`

	// +govalid:excluded_without_all=AdminUser AdminPassword
	GuestMode bool `json:"guest_mode"`

	DatabaseHost string `json:"database_host"`
	DatabasePort int    `json:"database_port"`

	// +govalid:excluded_without_all=DatabaseHost DatabasePort
	LocalStorageOnly bool `json:"local_storage_only"`
}
