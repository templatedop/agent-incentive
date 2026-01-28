//go:generate ./excluded_without.go
package excluded_without

// Feature is a struct for testing excluded_without validation
type Feature struct {
	PremiumAccess bool `json:"premium_access"`

	// +govalid:excluded_without=PremiumAccess
	AdvancedFeatures string `json:"advanced_features"`

	LicenseKey string `json:"license_key"`

	// +govalid:excluded_without=LicenseKey
	EnterpriseFeatures string `json:"enterprise_features"`
}
