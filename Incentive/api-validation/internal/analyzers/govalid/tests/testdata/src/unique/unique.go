//go:generate ./unique.go
package unique

// Collection is a struct for testing unique validation
type Collection struct {
	// +govalid:unique
	Tags []string `json:"tags"`

	// +govalid:unique
	IDs []int `json:"ids"`

	// +govalid:unique
	Categories []string `json:"categories"`
}
