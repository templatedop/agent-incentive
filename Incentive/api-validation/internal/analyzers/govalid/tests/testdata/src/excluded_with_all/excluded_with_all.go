//go:generate ./excluded_with_all.go
package excluded_with_all

// Config is a struct for testing excluded_with_all validation
type Config struct {
	CacheEnabled bool   `json:"cache_enabled"`
	CacheSize    int    `json:"cache_size"`

	// +govalid:excluded_with_all=CacheEnabled CacheSize
	DisableCache bool `json:"disable_cache"`

	SSLEnabled bool   `json:"ssl_enabled"`
	SSLCert    string `json:"ssl_cert"`

	// +govalid:excluded_with_all=SSLEnabled SSLCert
	InsecureMode bool `json:"insecure_mode"`
}
