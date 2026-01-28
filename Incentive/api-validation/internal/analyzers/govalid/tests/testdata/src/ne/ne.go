//go:generate ./ne.go
package ne

// Config is a struct for testing ne validation
type Config struct {
	// +govalid:ne=disabled
	Status string `json:"status"`

	// +govalid:ne=0
	Port int `json:"port"`

	// +govalid:ne=-1
	ErrorCode int `json:"error_code"`
}
