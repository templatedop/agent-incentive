package workflows

import (
	"fmt"

	config "gitlab.cept.gov.in/it-2.0-common/api-config"
	"go.temporal.io/sdk/client"
)

// NewTemporalClient creates a new Temporal client
func NewTemporalClient(cfg *config.Config) (client.Client, error) {
	// Get Temporal configuration from config
	host := cfg.GetString("temporal.host")
	port := cfg.GetString("temporal.port")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "7233"
	}
	hostPort := fmt.Sprintf("%s:%s", host, port)

	namespace := cfg.GetString("temporal.namespace")
	if namespace == "" {
		namespace = "default"
	}

	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Temporal client: %w", err)
	}

	return c, nil
}

// CloseTemporalClient closes the Temporal client
func CloseTemporalClient(c client.Client) {
	if c != nil {
		c.Close()
	}
}
