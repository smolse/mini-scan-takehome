package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// DataStoreConfig represents the configuration of the scan data store.
type DataStoreConfig struct {
	Type string `default:"cockroachdb"`

	// Cockroach configuration used when the data store type is "cockroachdb"
	CockroachHost     string `default:"cockroach"`
	CockroachPort     int    `default:"26257"`
	CockroachUser     string `default:"root"`
	CockroachDatabase string `default:"defaultdb"`
	CockroachSchema   string `default:"miniscan"`
	CockroachTable    string `default:"scans"`
}

// PubSubConfig represents the configuration of the Pub/Sub source.
type PubSubConfig struct {
	ProjectId              string `default:"test-project"`
	SubscriptionId         string `default:"scan-sub"`
	MaxOutstandingMessages int    `default:"10"`
}

// ServiceConfig represents the configuration of the service.
type ServiceConfig struct {
	GracefulShutdownTimeout time.Duration `default:"5s"`
}

// ProcessorConfig represents the entire configuration of the processor.
type ProcessorConfig struct {
	DataStore DataStoreConfig
	PubSub    PubSubConfig
	Service   ServiceConfig
}

// LoadProcessorConfig loads the processor configuration from the environment.
func LoadProcessorConfig() (ProcessorConfig, error) {
	var cfg ProcessorConfig
	err := envconfig.Process("PROCESSOR", &cfg)
	if err != nil {
		return ProcessorConfig{}, err
	}
	return cfg, nil
}
