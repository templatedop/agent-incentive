//go:generate ./maxduration.go
package maxduration

import "time"

// Request is a struct for testing maxduration validation
type Request struct {
	// +govalid:maxduration=10m
	Timeout time.Duration `json:"timeout"`

	// +govalid:maxduration=1h
	MaxWait time.Duration `json:"max_wait"`

	// +govalid:maxduration=30s
	Delay time.Duration `json:"delay"`
}
