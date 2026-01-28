//go:generate ./minduration.go
package minduration

import "time"

// Task is a struct for testing minduration validation
type Task struct {
	// +govalid:minduration=1h
	Duration time.Duration `json:"duration"`

	// +govalid:minduration=30s
	Timeout time.Duration `json:"timeout"`

	// +govalid:minduration=5m
	Interval time.Duration `json:"interval"`
}
