//go:generate ./isdefault.go
package isdefault

// Optional is a struct for testing isdefault validation
type Optional struct {
	// +govalid:isdefault
	EmptyString string `json:"empty_string"`

	// +govalid:isdefault
	ZeroInt int `json:"zero_int"`

	// +govalid:isdefault
	FalseBool bool `json:"false_bool"`
}
